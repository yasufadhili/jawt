package lsp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"slices"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/yasufadhili/jawt/internal/tsc/collections"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/ls"
	"github.com/yasufadhili/jawt/internal/tsc/lsp/lsproto"
	"github.com/yasufadhili/jawt/internal/tsc/project"
	"github.com/yasufadhili/jawt/internal/tsc/vfs"
	"golang.org/x/sync/errgroup"
)

type ServerOptions struct {
	In  Reader
	Out Writer
	Err io.Writer

	Cwd                string
	NewLine            core.NewLineKind
	FS                 vfs.FS
	DefaultLibraryPath string
	TypingsLocation    string

	ParsedFileCache project.ParsedFileCache
}

func NewServer(opts *ServerOptions) *Server {
	if opts.Cwd == "" {
		panic("Cwd is required")
	}
	return &Server{
		r:                     opts.In,
		w:                     opts.Out,
		stderr:                opts.Err,
		requestQueue:          make(chan *lsproto.RequestMessage, 100),
		outgoingQueue:         make(chan *lsproto.Message, 100),
		pendingClientRequests: make(map[lsproto.ID]pendingClientRequest),
		pendingServerRequests: make(map[lsproto.ID]chan *lsproto.ResponseMessage),
		cwd:                   opts.Cwd,
		newLine:               opts.NewLine,
		fs:                    opts.FS,
		defaultLibraryPath:    opts.DefaultLibraryPath,
		typingsLocation:       opts.TypingsLocation,
		parsedFileCache:       opts.ParsedFileCache,
	}
}

var (
	_ project.ServiceHost = (*Server)(nil)
	_ project.Client      = (*Server)(nil)
)

type pendingClientRequest struct {
	req    *lsproto.RequestMessage
	cancel context.CancelFunc
}

type Reader interface {
	Read() (*lsproto.Message, error)
}

type Writer interface {
	Write(msg *lsproto.Message) error
}

type lspReader struct {
	r *lsproto.BaseReader
}

type lspWriter struct {
	w *lsproto.BaseWriter
}

func (r *lspReader) Read() (*lsproto.Message, error) {
	data, err := r.r.Read()
	if err != nil {
		return nil, err
	}

	req := &lsproto.Message{}
	if err := json.Unmarshal(data, req); err != nil {
		return nil, fmt.Errorf("%w: %w", lsproto.ErrInvalidRequest, err)
	}

	return req, nil
}

func ToReader(r io.Reader) Reader {
	return &lspReader{r: lsproto.NewBaseReader(r)}
}

func (w *lspWriter) Write(msg *lsproto.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	return w.w.Write(data)
}

func ToWriter(w io.Writer) Writer {
	return &lspWriter{w: lsproto.NewBaseWriter(w)}
}

var (
	_ Reader = (*lspReader)(nil)
	_ Writer = (*lspWriter)(nil)
)

type Server struct {
	r Reader
	w Writer

	stderr io.Writer

	clientSeq               atomic.Int32
	requestQueue            chan *lsproto.RequestMessage
	outgoingQueue           chan *lsproto.Message
	pendingClientRequests   map[lsproto.ID]pendingClientRequest
	pendingClientRequestsMu sync.Mutex
	pendingServerRequests   map[lsproto.ID]chan *lsproto.ResponseMessage
	pendingServerRequestsMu sync.Mutex

	cwd                string
	newLine            core.NewLineKind
	fs                 vfs.FS
	defaultLibraryPath string
	typingsLocation    string

	initializeParams *lsproto.InitializeParams
	positionEncoding lsproto.PositionEncodingKind

	watchEnabled bool
	watcherID    atomic.Uint32
	watchers     collections.SyncSet[project.WatcherHandle]

	logger         *project.Logger
	projectService *project.Service

	// enables tests to share a cache of parsed source files
	parsedFileCache project.ParsedFileCache

	// !!! temporary; remove when we have `handleDidChangeConfiguration`/implicit project config support
	compilerOptionsForInferredProjects *core.CompilerOptions
}

// FS implements project.ServiceHost.
func (s *Server) FS() vfs.FS {
	return s.fs
}

// DefaultLibraryPath implements project.ServiceHost.
func (s *Server) DefaultLibraryPath() string {
	return s.defaultLibraryPath
}

// TypingsLocation implements project.ServiceHost.
func (s *Server) TypingsLocation() string {
	return s.typingsLocation
}

// GetCurrentDirectory implements project.ServiceHost.
func (s *Server) GetCurrentDirectory() string {
	return s.cwd
}

// NewLine implements project.ServiceHost.
func (s *Server) NewLine() string {
	return s.newLine.GetNewLineCharacter()
}

// Trace implements project.ServiceHost.
func (s *Server) Trace(msg string) {
	s.Log(msg)
}

// Client implements project.ServiceHost.
func (s *Server) Client() project.Client {
	if !s.watchEnabled {
		return nil
	}
	return s
}

// WatchFiles implements project.Client.
func (s *Server) WatchFiles(ctx context.Context, watchers []*lsproto.FileSystemWatcher) (project.WatcherHandle, error) {
	watcherId := fmt.Sprintf("watcher-%d", s.watcherID.Add(1))
	_, err := s.sendRequest(ctx, lsproto.MethodClientRegisterCapability, &lsproto.RegistrationParams{
		Registrations: []*lsproto.Registration{
			{
				Id:     watcherId,
				Method: string(lsproto.MethodWorkspaceDidChangeWatchedFiles),
				RegisterOptions: ptrTo(any(lsproto.DidChangeWatchedFilesRegistrationOptions{
					Watchers: watchers,
				})),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to register file watcher: %w", err)
	}

	handle := project.WatcherHandle(watcherId)
	s.watchers.Add(handle)
	return handle, nil
}

// UnwatchFiles implements project.Client.
func (s *Server) UnwatchFiles(ctx context.Context, handle project.WatcherHandle) error {
	if s.watchers.Has(handle) {
		_, err := s.sendRequest(ctx, lsproto.MethodClientUnregisterCapability, &lsproto.UnregistrationParams{
			Unregisterations: []*lsproto.Unregistration{
				{
					Id:     string(handle),
					Method: string(lsproto.MethodWorkspaceDidChangeWatchedFiles),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to unregister file watcher: %w", err)
		}

		s.watchers.Delete(handle)
		return nil
	}

	return fmt.Errorf("no file watcher exists with ID %s", handle)
}

// RefreshDiagnostics implements project.Client.
func (s *Server) RefreshDiagnostics(ctx context.Context) error {
	if ptrIsTrue(s.initializeParams.Capabilities.Workspace.Diagnostics.RefreshSupport) {
		if _, err := s.sendRequest(ctx, lsproto.MethodWorkspaceDiagnosticRefresh, nil); err != nil {
			return fmt.Errorf("failed to refresh diagnostics: %w", err)
		}
	}
	return nil
}

func (s *Server) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return s.dispatchLoop(ctx) })
	g.Go(func() error { return s.writeLoop(ctx) })

	// Don't run readLoop in the group, as it blocks on stdin read and cannot be cancelled.
	readLoopErr := make(chan error, 1)
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-readLoopErr:
			return err
		}
	})
	go func() { readLoopErr <- s.readLoop(ctx) }()

	if err := g.Wait(); err != nil && !errors.Is(err, io.EOF) && ctx.Err() != nil {
		return err
	}
	return nil
}

func (s *Server) readLoop(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		msg, err := s.read()
		if err != nil {
			if errors.Is(err, lsproto.ErrInvalidRequest) {
				s.sendError(nil, err)
				continue
			}
			return err
		}

		if s.initializeParams == nil && msg.Kind == lsproto.MessageKindRequest {
			req := msg.AsRequest()
			if req.Method == lsproto.MethodInitialize {
				s.handleInitialize(req)
			} else {
				s.sendError(req.ID, lsproto.ErrServerNotInitialized)
			}
			continue
		}

		if msg.Kind == lsproto.MessageKindResponse {
			resp := msg.AsResponse()
			s.pendingServerRequestsMu.Lock()
			if respChan, ok := s.pendingServerRequests[*resp.ID]; ok {
				respChan <- resp
				close(respChan)
				delete(s.pendingServerRequests, *resp.ID)
			}
			s.pendingServerRequestsMu.Unlock()
		} else {
			req := msg.AsRequest()
			if req.Method == lsproto.MethodCancelRequest {
				s.cancelRequest(req.Params.(*lsproto.CancelParams).Id)
			} else {
				s.requestQueue <- req
			}
		}
	}
}

func (s *Server) cancelRequest(rawID lsproto.IntegerOrString) {
	id := lsproto.NewID(rawID)
	s.pendingClientRequestsMu.Lock()
	defer s.pendingClientRequestsMu.Unlock()
	if pendingReq, ok := s.pendingClientRequests[*id]; ok {
		pendingReq.cancel()
		delete(s.pendingClientRequests, *id)
	}
}

func (s *Server) read() (*lsproto.Message, error) {
	return s.r.Read()
}

func (s *Server) dispatchLoop(ctx context.Context) error {
	ctx, lspExit := context.WithCancel(ctx)
	defer lspExit()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case req := <-s.requestQueue:
			requestCtx := ctx
			if req.ID != nil {
				var cancel context.CancelFunc
				requestCtx, cancel = context.WithCancel(core.WithRequestID(requestCtx, req.ID.String()))
				s.pendingClientRequestsMu.Lock()
				s.pendingClientRequests[*req.ID] = pendingClientRequest{
					req:    req,
					cancel: cancel,
				}
				s.pendingClientRequestsMu.Unlock()
			}

			handle := func() {
				defer func() {
					if r := recover(); r != nil {
						stack := debug.Stack()
						s.Log("panic handling request", req.Method, r, string(stack))
						if isBlockingMethod(req.Method) {
							lspExit()
						} else {
							if req.ID != nil {
								s.sendError(req.ID, fmt.Errorf("%w: panic handling request %s: %v", lsproto.ErrInternalError, req.Method, r))
							} else {
								s.Log("unhandled panic in notification", req.Method, r)
							}
						}
					}
				}()
				if err := s.handleRequestOrNotification(requestCtx, req); err != nil {
					if errors.Is(err, io.EOF) {
						lspExit()
					} else {
						s.sendError(req.ID, err)
					}
				}

				if req.ID != nil {
					s.pendingClientRequestsMu.Lock()
					delete(s.pendingClientRequests, *req.ID)
					s.pendingClientRequestsMu.Unlock()
				}
			}

			if isBlockingMethod(req.Method) {
				handle()
			} else {
				go handle()
			}
		}
	}
}

func (s *Server) writeLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-s.outgoingQueue:
			if err := s.w.Write(msg); err != nil {
				return fmt.Errorf("failed to write message: %w", err)
			}
		}
	}
}

func (s *Server) sendRequest(ctx context.Context, method lsproto.Method, params any) (any, error) {
	id := lsproto.NewIDString(fmt.Sprintf("ts%d", s.clientSeq.Add(1)))
	req := lsproto.NewRequestMessage(method, id, params)

	responseChan := make(chan *lsproto.ResponseMessage, 1)
	s.pendingServerRequestsMu.Lock()
	s.pendingServerRequests[*id] = responseChan
	s.pendingServerRequestsMu.Unlock()

	s.outgoingQueue <- req.Message()

	select {
	case <-ctx.Done():
		s.pendingServerRequestsMu.Lock()
		defer s.pendingServerRequestsMu.Unlock()
		if respChan, ok := s.pendingServerRequests[*id]; ok {
			close(respChan)
			delete(s.pendingServerRequests, *id)
		}
		return nil, ctx.Err()
	case resp := <-responseChan:
		if resp.Error != nil {
			return nil, fmt.Errorf("request failed: %s", resp.Error.String())
		}
		return resp.Result, nil
	}
}

func (s *Server) sendResult(id *lsproto.ID, result any) {
	s.sendResponse(&lsproto.ResponseMessage{
		ID:     id,
		Result: result,
	})
}

func (s *Server) sendError(id *lsproto.ID, err error) {
	code := lsproto.ErrInternalError.Code
	if errCode := (*lsproto.ErrorCode)(nil); errors.As(err, &errCode) {
		code = errCode.Code
	}
	// TODO(jakebailey): error data
	s.sendResponse(&lsproto.ResponseMessage{
		ID: id,
		Error: &lsproto.ResponseError{
			Code:    code,
			Message: err.Error(),
		},
	})
}

func (s *Server) sendResponse(resp *lsproto.ResponseMessage) {
	s.outgoingQueue <- resp.Message()
}

func (s *Server) handleRequestOrNotification(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params
	switch params.(type) {
	case *lsproto.InitializeParams:
		s.sendError(req.ID, lsproto.ErrInvalidRequest)
		return nil
	case *lsproto.InitializedParams:
		return s.handleInitialized(ctx, req)
	case *lsproto.DidOpenTextDocumentParams:
		return s.handleDidOpen(ctx, req)
	case *lsproto.DidChangeTextDocumentParams:
		return s.handleDidChange(ctx, req)
	case *lsproto.DidSaveTextDocumentParams:
		return s.handleDidSave(ctx, req)
	case *lsproto.DidCloseTextDocumentParams:
		return s.handleDidClose(ctx, req)
	case *lsproto.DidChangeWatchedFilesParams:
		return s.handleDidChangeWatchedFiles(ctx, req)
	case *lsproto.DocumentDiagnosticParams:
		return s.handleDocumentDiagnostic(ctx, req)
	case *lsproto.HoverParams:
		return s.handleHover(ctx, req)
	case *lsproto.DefinitionParams:
		return s.handleDefinition(ctx, req)
	case *lsproto.CompletionParams:
		return s.handleCompletion(ctx, req)
	case *lsproto.ReferenceParams:
		return s.handleReferences(ctx, req)
	case *lsproto.SignatureHelpParams:
		return s.handleSignatureHelp(ctx, req)
	case *lsproto.DocumentFormattingParams:
		return s.handleDocumentFormat(ctx, req)
	case *lsproto.DocumentRangeFormattingParams:
		return s.handleDocumentRangeFormat(ctx, req)
	case *lsproto.DocumentOnTypeFormattingParams:
		return s.handleDocumentOnTypeFormat(ctx, req)
	case *lsproto.WorkspaceSymbolParams:
		return s.handleWorkspaceSymbol(ctx, req)
	default:
		switch req.Method {
		case lsproto.MethodShutdown:
			s.projectService.Close()
			s.sendResult(req.ID, nil)
			return nil
		case lsproto.MethodExit:
			return io.EOF
		default:
			s.Log("unknown method", req.Method)
			if req.ID != nil {
				s.sendError(req.ID, lsproto.ErrInvalidRequest)
			}
			return nil
		}
	}
}

func (s *Server) handleInitialize(req *lsproto.RequestMessage) {
	s.initializeParams = req.Params.(*lsproto.InitializeParams)

	s.positionEncoding = lsproto.PositionEncodingKindUTF16
	if genCapabilities := s.initializeParams.Capabilities.General; genCapabilities != nil && genCapabilities.PositionEncodings != nil {
		if slices.Contains(*genCapabilities.PositionEncodings, lsproto.PositionEncodingKindUTF8) {
			s.positionEncoding = lsproto.PositionEncodingKindUTF8
		}
	}

	s.sendResult(req.ID, &lsproto.InitializeResult{
		ServerInfo: &lsproto.ServerInfo{
			Name:    "typescript-go",
			Version: ptrTo(core.Version()),
		},
		Capabilities: &lsproto.ServerCapabilities{
			PositionEncoding: ptrTo(s.positionEncoding),
			TextDocumentSync: &lsproto.TextDocumentSyncOptionsOrTextDocumentSyncKind{
				TextDocumentSyncOptions: &lsproto.TextDocumentSyncOptions{
					OpenClose: ptrTo(true),
					Change:    ptrTo(lsproto.TextDocumentSyncKindIncremental),
					Save: &lsproto.BooleanOrSaveOptions{
						SaveOptions: &lsproto.SaveOptions{
							IncludeText: ptrTo(true),
						},
					},
				},
			},
			HoverProvider: &lsproto.BooleanOrHoverOptions{
				Boolean: ptrTo(true),
			},
			DefinitionProvider: &lsproto.BooleanOrDefinitionOptions{
				Boolean: ptrTo(true),
			},
			ReferencesProvider: &lsproto.BooleanOrReferenceOptions{
				Boolean: ptrTo(true),
			},
			DiagnosticProvider: &lsproto.DiagnosticOptionsOrDiagnosticRegistrationOptions{
				DiagnosticOptions: &lsproto.DiagnosticOptions{
					InterFileDependencies: true,
				},
			},
			CompletionProvider: &lsproto.CompletionOptions{
				TriggerCharacters: &ls.TriggerCharacters,
				// !!! other options
			},
			SignatureHelpProvider: &lsproto.SignatureHelpOptions{
				TriggerCharacters: &[]string{"(", ","},
			},
			DocumentFormattingProvider: &lsproto.BooleanOrDocumentFormattingOptions{
				Boolean: ptrTo(true),
			},
			DocumentRangeFormattingProvider: &lsproto.BooleanOrDocumentRangeFormattingOptions{
				Boolean: ptrTo(true),
			},
			DocumentOnTypeFormattingProvider: &lsproto.DocumentOnTypeFormattingOptions{
				FirstTriggerCharacter: "{",
				MoreTriggerCharacter:  &[]string{"}", ";", "\n"},
			},
			WorkspaceSymbolProvider: &lsproto.BooleanOrWorkspaceSymbolOptions{
				Boolean: ptrTo(true),
			},
		},
	})
}

func (s *Server) handleInitialized(ctx context.Context, req *lsproto.RequestMessage) error {
	if shouldEnableWatch(s.initializeParams) {
		s.watchEnabled = true
	}

	s.logger = project.NewLogger([]io.Writer{s.stderr}, "" /*file*/, project.LogLevelVerbose)
	s.projectService = project.NewService(s, project.ServiceOptions{
		Logger:           s.logger,
		WatchEnabled:     s.watchEnabled,
		PositionEncoding: s.positionEncoding,
		TypingsInstallerOptions: project.TypingsInstallerOptions{
			ThrottleLimit: 5,
			NpmInstall:    project.NpmInstall,
		},
		ParsedFileCache: s.parsedFileCache,
	})
	// !!! temporary; remove when we have `handleDidChangeConfiguration`/implicit project config support
	if s.compilerOptionsForInferredProjects != nil {
		s.projectService.SetCompilerOptionsForInferredProjects(s.compilerOptionsForInferredProjects)
	}

	return nil
}

func (s *Server) handleDidOpen(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidOpenTextDocumentParams)
	s.projectService.OpenFile(ls.DocumentURIToFileName(params.TextDocument.Uri), params.TextDocument.Text, ls.LanguageKindToScriptKind(params.TextDocument.LanguageId), "")
	return nil
}

func (s *Server) handleDidChange(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidChangeTextDocumentParams)
	return s.projectService.ChangeFile(params.TextDocument, params.ContentChanges)
}

func (s *Server) handleDidSave(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidSaveTextDocumentParams)
	s.projectService.MarkFileSaved(ls.DocumentURIToFileName(params.TextDocument.Uri), *params.Text)
	return nil
}

func (s *Server) handleDidClose(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidCloseTextDocumentParams)
	s.projectService.CloseFile(ls.DocumentURIToFileName(params.TextDocument.Uri))
	return nil
}

func (s *Server) handleDidChangeWatchedFiles(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidChangeWatchedFilesParams)
	return s.projectService.OnWatchedFilesChanged(ctx, params.Changes)
}

func (s *Server) handleDocumentDiagnostic(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DocumentDiagnosticParams)
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	diagnostics, err := languageService.GetDocumentDiagnostics(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, diagnostics)
	return nil
}

func (s *Server) handleHover(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.HoverParams)
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	hover, err := languageService.ProvideHover(ctx, params.TextDocument.Uri, params.Position)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, hover)
	return nil
}

func (s *Server) handleSignatureHelp(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.SignatureHelpParams)
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	signatureHelp := languageService.ProvideSignatureHelp(
		ctx,
		params.TextDocument.Uri,
		params.Position,
		params.Context,
		s.initializeParams.Capabilities.TextDocument.SignatureHelp,
		&ls.UserPreferences{},
	)
	s.sendResult(req.ID, signatureHelp)
	return nil
}

func (s *Server) handleDefinition(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DefinitionParams)
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	definition, err := languageService.ProvideDefinition(ctx, params.TextDocument.Uri, params.Position)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, definition)
	return nil
}

func (s *Server) handleReferences(ctx context.Context, req *lsproto.RequestMessage) error {
	// findAllReferences
	params := req.Params.(*lsproto.ReferenceParams)
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	locations := languageService.ProvideReferences(params)
	s.sendResult(req.ID, locations)
	return nil
}

func (s *Server) handleCompletion(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.CompletionParams)
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	// !!! get user preferences
	list, err := languageService.ProvideCompletion(
		ctx,
		params.TextDocument.Uri,
		params.Position,
		params.Context,
		getCompletionClientCapabilities(s.initializeParams),
		&ls.UserPreferences{})
	if err != nil {
		return err
	}
	s.sendResult(req.ID, list)
	return nil
}

func (s *Server) handleDocumentFormat(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DocumentFormattingParams)
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	res, err := languageService.ProvideFormatDocument(
		ctx,
		params.TextDocument.Uri,
		params.Options,
	)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, res)
	return nil
}

func (s *Server) handleDocumentRangeFormat(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DocumentRangeFormattingParams)
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	res, err := languageService.ProvideFormatDocumentRange(
		ctx,
		params.TextDocument.Uri,
		params.Options,
		params.Range,
	)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, res)
	return nil
}

func (s *Server) handleDocumentOnTypeFormat(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DocumentOnTypeFormattingParams)
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	res, err := languageService.ProvideFormatDocumentOnType(
		ctx,
		params.TextDocument.Uri,
		params.Options,
		params.Position,
		params.Ch,
	)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, res)
	return nil
}

func (s *Server) handleWorkspaceSymbol(ctx context.Context, req *lsproto.RequestMessage) error {
	programs := core.Map(s.projectService.Projects(), (*project.Project).GetProgram)
	params := req.Params.(*lsproto.WorkspaceSymbolParams)
	symbols, err := ls.ProvideWorkspaceSymbols(ctx, programs, s.projectService.Converters(), params.Query)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, symbols)
	return nil
}

func (s *Server) Log(msg ...any) {
	fmt.Fprintln(s.stderr, msg...)
}

// !!! temporary; remove when we have `handleDidChangeConfiguration`/implicit project config support
func (s *Server) SetCompilerOptionsForInferredProjects(options *core.CompilerOptions) {
	s.compilerOptionsForInferredProjects = options
	if s.projectService != nil {
		s.projectService.SetCompilerOptionsForInferredProjects(options)
	}
}

func isBlockingMethod(method lsproto.Method) bool {
	switch method {
	case lsproto.MethodInitialize,
		lsproto.MethodInitialized,
		lsproto.MethodTextDocumentDidOpen,
		lsproto.MethodTextDocumentDidChange,
		lsproto.MethodTextDocumentDidSave,
		lsproto.MethodTextDocumentDidClose,
		lsproto.MethodWorkspaceDidChangeWatchedFiles:
		return true
	}
	return false
}

func ptrTo[T any](v T) *T {
	return &v
}

func ptrIsTrue(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func shouldEnableWatch(params *lsproto.InitializeParams) bool {
	if params == nil || params.Capabilities == nil || params.Capabilities.Workspace == nil {
		return false
	}
	return params.Capabilities.Workspace.DidChangeWatchedFiles != nil &&
		ptrIsTrue(params.Capabilities.Workspace.DidChangeWatchedFiles.DynamicRegistration)
}

func getCompletionClientCapabilities(params *lsproto.InitializeParams) *lsproto.CompletionClientCapabilities {
	if params == nil || params.Capabilities == nil || params.Capabilities.TextDocument == nil {
		return nil
	}
	return params.Capabilities.TextDocument.Completion
}
