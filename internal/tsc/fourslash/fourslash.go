package fourslash

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/bundled"
	"github.com/yasufadhili/jawt/internal/tsc/collections"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/ls"
	"github.com/yasufadhili/jawt/internal/tsc/lsp"
	"github.com/yasufadhili/jawt/internal/tsc/lsp/lsproto"
	"github.com/yasufadhili/jawt/internal/tsc/project"
	"github.com/yasufadhili/jawt/internal/tsc/testutil/harnessutil"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
	"github.com/yasufadhili/jawt/internal/tsc/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type FourslashTest struct {
	server *lsp.Server
	in     *lspWriter
	out    *lspReader
	id     int32

	testData *TestData

	currentCaretPosition lsproto.Position
	currentFilename      string
	lastKnownMarkerName  string
	activeFilename       string
}

type lspReader struct {
	c <-chan *lsproto.Message
}

func (r *lspReader) Read() (*lsproto.Message, error) {
	msg, ok := <-r.c
	if !ok {
		return nil, io.EOF
	}
	return msg, nil
}

type lspWriter struct {
	c chan<- *lsproto.Message
}

func (w *lspWriter) Write(msg *lsproto.Message) error {
	w.c <- msg
	return nil
}

func (r *lspWriter) Close() {
	close(r.c)
}

var (
	_ lsp.Reader = (*lspReader)(nil)
	_ lsp.Writer = (*lspWriter)(nil)
)

func newLSPPipe() (*lspReader, *lspWriter) {
	c := make(chan *lsproto.Message, 100)
	return &lspReader{c: c}, &lspWriter{c: c}
}

var sourceFileCache collections.SyncMap[harnessutil.SourceFileCacheKey, *ast.SourceFile]

type parsedFileCache struct{}

func (c *parsedFileCache) GetFile(opts ast.SourceFileParseOptions, text string, scriptKind core.ScriptKind) *ast.SourceFile {
	key := harnessutil.GetSourceFileCacheKey(opts, text, scriptKind)
	cachedFile, ok := sourceFileCache.Load(key)
	if !ok {
		return nil
	}
	return cachedFile
}

func (c *parsedFileCache) CacheFile(opts ast.SourceFileParseOptions, text string, scriptKind core.ScriptKind, sourceFile *ast.SourceFile) {
	key := harnessutil.GetSourceFileCacheKey(opts, text, scriptKind)
	sourceFileCache.Store(key, sourceFile)
}

var _ project.ParsedFileCache = (*parsedFileCache)(nil)

func NewFourslash(t *testing.T, capabilities *lsproto.ClientCapabilities, content string) *FourslashTest {
	if !bundled.Embedded {
		// Without embedding, we'd need to read all of the lib files out from disk into the MapFS.
		// Just skip this for now.
		t.Skip("bundled files are not embedded")
	}
	rootDir := "/"
	fileName := getFileNameFromTest(t)
	testfs := make(map[string]string)
	testData := ParseTestData(t, content, fileName)
	for _, file := range testData.Files {
		filePath := tspath.GetNormalizedAbsolutePath(file.fileName, rootDir)
		testfs[filePath] = file.Content
	}

	compilerOptions := &core.CompilerOptions{}
	harnessutil.SetCompilerOptionsFromTestConfig(t, testData.GlobalOptions, compilerOptions)
	compilerOptions.SkipDefaultLibCheck = core.TSTrue

	inputReader, inputWriter := newLSPPipe()
	outputReader, outputWriter := newLSPPipe()
	fs := vfstest.FromMap(testfs, true /*useCaseSensitiveFileNames*/)

	var err strings.Builder
	server := lsp.NewServer(&lsp.ServerOptions{
		In:  inputReader,
		Out: outputWriter,
		Err: &err,

		Cwd:                "/",
		NewLine:            core.NewLineKindLF,
		FS:                 bundled.WrapFS(fs),
		DefaultLibraryPath: bundled.LibPath(),

		ParsedFileCache: &parsedFileCache{},
	})

	go func() {
		defer func() {
			outputWriter.Close()
		}()
		err := server.Run()
		if err != nil {
			t.Error("server error:", err)
		}
	}()

	f := &FourslashTest{
		server:   server,
		in:       inputWriter,
		out:      outputReader,
		testData: &testData,
	}

	// !!! temporary; remove when we have `handleDidChangeConfiguration`/implicit project config support
	// !!! replace with a proper request *after initialize*
	f.server.SetCompilerOptionsForInferredProjects(compilerOptions)
	f.initialize(t, capabilities)
	f.openFile(t, f.testData.Files[0])

	t.Cleanup(func() {
		inputWriter.Close()
	})
	return f
}

func getFileNameFromTest(t *testing.T) string {
	name := strings.TrimPrefix(t.Name(), "Test")
	char, size := utf8.DecodeRuneInString(name)
	return string(unicode.ToLower(char)) + name[size:] + tspath.ExtensionTs
}

func (f *FourslashTest) nextID() int32 {
	id := f.id
	f.id++
	return id
}

func (f *FourslashTest) initialize(t *testing.T, capabilities *lsproto.ClientCapabilities) {
	params := &lsproto.InitializeParams{}
	params.Capabilities = getCapabilitiesWithDefaults(capabilities)
	// !!! check for errors?
	f.sendRequest(t, lsproto.MethodInitialize, params)
	f.sendNotification(t, lsproto.MethodInitialized, &lsproto.InitializedParams{})
}

var (
	ptrTrue                       = ptrTo(true)
	defaultCompletionCapabilities = &lsproto.CompletionClientCapabilities{
		CompletionItem: &lsproto.ClientCompletionItemOptions{
			SnippetSupport:          ptrTrue,
			CommitCharactersSupport: ptrTrue,
			PreselectSupport:        ptrTrue,
			LabelDetailsSupport:     ptrTrue,
			InsertReplaceSupport:    ptrTrue,
		},
		CompletionList: &lsproto.CompletionListCapabilities{
			ItemDefaults: &[]string{"commitCharacters", "editRange"},
		},
	}
)

func getCapabilitiesWithDefaults(capabilities *lsproto.ClientCapabilities) *lsproto.ClientCapabilities {
	var capabilitiesWithDefaults lsproto.ClientCapabilities
	if capabilities != nil {
		capabilitiesWithDefaults = *capabilities
	}
	capabilitiesWithDefaults.General = &lsproto.GeneralClientCapabilities{
		PositionEncodings: &[]lsproto.PositionEncodingKind{lsproto.PositionEncodingKindUTF8},
	}
	if capabilitiesWithDefaults.TextDocument == nil {
		capabilitiesWithDefaults.TextDocument = &lsproto.TextDocumentClientCapabilities{}
	}
	if capabilitiesWithDefaults.TextDocument.Completion == nil {
		capabilitiesWithDefaults.TextDocument.Completion = defaultCompletionCapabilities
	}
	return &capabilitiesWithDefaults
}

func (f *FourslashTest) sendRequest(t *testing.T, method lsproto.Method, params any) *lsproto.Message {
	id := f.nextID()
	req := lsproto.NewRequestMessage(
		method,
		lsproto.NewID(lsproto.IntegerOrString{Integer: &id}),
		params,
	)
	f.writeMsg(t, req.Message())
	return f.readMsg(t)
}

func (f *FourslashTest) sendNotification(t *testing.T, method lsproto.Method, params any) {
	notification := lsproto.NewNotificationMessage(
		method,
		params,
	)
	f.writeMsg(t, notification.Message())
}

func (f *FourslashTest) writeMsg(t *testing.T, msg *lsproto.Message) {
	if err := f.in.Write(msg); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}
}

func (f *FourslashTest) readMsg(t *testing.T) *lsproto.Message {
	// !!! filter out response by id etc
	msg, err := f.out.Read()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}
	return msg
}

func (f *FourslashTest) GoToMarker(t *testing.T, markerName string) {
	marker, ok := f.testData.MarkerPositions[markerName]
	if !ok {
		t.Fatalf("Marker %s not found", markerName)
	}
	f.ensureActiveFile(t, marker.FileName)
	f.currentCaretPosition = marker.LSPosition
	f.currentFilename = marker.FileName
	f.lastKnownMarkerName = marker.Name
}

func (f *FourslashTest) Markers() []*Marker {
	return f.testData.Markers
}

func (f *FourslashTest) Ranges() []*RangeMarker {
	return f.testData.Ranges
}

func (f *FourslashTest) ensureActiveFile(t *testing.T, filename string) {
	if f.activeFilename != filename {
		file := core.Find(f.testData.Files, func(f *TestFileInfo) bool {
			return f.fileName == filename
		})
		if file == nil {
			t.Fatalf("File %s not found in test data", filename)
		}
		f.openFile(t, file)
	}
}

func (f *FourslashTest) openFile(t *testing.T, file *TestFileInfo) {
	f.activeFilename = file.fileName
	f.sendNotification(t, lsproto.MethodTextDocumentDidOpen, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{
			Uri:        ls.FileNameToDocumentURI(file.fileName),
			LanguageId: getLanguageKind(file.fileName),
			Text:       file.Content,
		},
	})
}

func getLanguageKind(filename string) lsproto.LanguageKind {
	if tspath.FileExtensionIsOneOf(
		filename,
		[]string{
			tspath.ExtensionTs, tspath.ExtensionMts, tspath.ExtensionCts,
			tspath.ExtensionDmts, tspath.ExtensionDcts, tspath.ExtensionDts,
		}) {
		return lsproto.LanguageKindTypeScript
	}
	if tspath.FileExtensionIsOneOf(filename, []string{tspath.ExtensionJs, tspath.ExtensionMjs, tspath.ExtensionCjs}) {
		return lsproto.LanguageKindJavaScript
	}
	if tspath.FileExtensionIs(filename, tspath.ExtensionJsx) {
		return lsproto.LanguageKindJavaScriptReact
	}
	if tspath.FileExtensionIs(filename, tspath.ExtensionTsx) {
		return lsproto.LanguageKindTypeScriptReact
	}
	if tspath.FileExtensionIs(filename, tspath.ExtensionJson) {
		return lsproto.LanguageKindJSON
	}
	return lsproto.LanguageKindTypeScript // !!! should we error in this case?
}

type CompletionsExpectedList struct {
	IsIncomplete bool
	ItemDefaults *CompletionsExpectedItemDefaults
	Items        *CompletionsExpectedItems
}

type Ignored = struct{}

// *EditRange | Ignored
type ExpectedCompletionEditRange = any

type EditRange struct {
	Insert  *RangeMarker
	Replace *RangeMarker
}

type CompletionsExpectedItemDefaults struct {
	CommitCharacters *[]string
	EditRange        ExpectedCompletionEditRange
}

// *lsproto.CompletionItem | string
type CompletionsExpectedItem = any

// !!! unsorted completions
type CompletionsExpectedItems struct {
	Includes []CompletionsExpectedItem
	Excludes []string
	Exact    []CompletionsExpectedItem
}

// string | *Marker | []string | []*Marker
type MarkerInput = any

// !!! user preferences param
// !!! completion context param
// !!! go to marker: use current marker if none specified/support nil marker input
func (f *FourslashTest) VerifyCompletions(t *testing.T, markerInput MarkerInput, expected *CompletionsExpectedList) {
	switch marker := markerInput.(type) {
	case string:
		f.verifyCompletionsAtMarker(t, marker, expected)
	case *Marker:
		f.verifyCompletionsAtMarker(t, marker.Name, expected)
	case []string:
		for _, markerName := range marker {
			f.verifyCompletionsAtMarker(t, markerName, expected)
		}
	case []*Marker:
		for _, marker := range marker {
			f.verifyCompletionsAtMarker(t, marker.Name, expected)
		}
	case nil:
		f.verifyCompletionsWorker(t, expected)
	default:
		t.Fatalf("Invalid marker input type: %T. Expected string, *Marker, []string, or []*Marker.", markerInput)
	}
}

func (f *FourslashTest) verifyCompletionsAtMarker(t *testing.T, markerName string, expected *CompletionsExpectedList) {
	f.GoToMarker(t, markerName)
	f.verifyCompletionsWorker(t, expected)
}

func (f *FourslashTest) verifyCompletionsWorker(t *testing.T, expected *CompletionsExpectedList) {
	params := &lsproto.CompletionParams{
		TextDocumentPositionParams: lsproto.TextDocumentPositionParams{
			TextDocument: lsproto.TextDocumentIdentifier{
				Uri: ls.FileNameToDocumentURI(f.currentFilename),
			},
			Position: f.currentCaretPosition,
		},
		Context: &lsproto.CompletionContext{},
	}
	resMsg := f.sendRequest(t, lsproto.MethodTextDocumentCompletion, params)
	if resMsg == nil {
		t.Fatalf("Nil response received for completion request at marker %s", f.lastKnownMarkerName)
	}
	result := resMsg.AsResponse().Result
	switch result := result.(type) {
	case *lsproto.CompletionList:
		verifyCompletionsResult(t, f.lastKnownMarkerName, result, expected)
	default:
		t.Fatalf("Unexpected response type for completion request at marker %s: %v", f.lastKnownMarkerName, result)
	}
}

func verifyCompletionsResult(t *testing.T, markerName string, actual *lsproto.CompletionList, expected *CompletionsExpectedList) {
	prefix := fmt.Sprintf("At marker '%s': ", markerName)
	if actual == nil {
		if !isEmptyExpectedList(expected) {
			t.Fatal(prefix + "Expected completion list but got nil.")
		}
		return
	} else if expected == nil {
		// !!! cmp.Diff(actual, nil) should probably be a .String() call here and elswhere
		t.Fatalf(prefix+"Expected nil completion list but got non-nil: %s", cmp.Diff(actual, nil))
	}
	assert.Equal(t, actual.IsIncomplete, expected.IsIncomplete, prefix+"IsIncomplete mismatch")
	verifyCompletionsItemDefaults(t, actual.ItemDefaults, expected.ItemDefaults, prefix+"ItemDefaults mismatch: ")
	verifyCompletionsItems(t, prefix, actual.Items, expected.Items)
}

func isEmptyExpectedList(expected *CompletionsExpectedList) bool {
	return expected == nil || (len(expected.Items.Exact) == 0 && len(expected.Items.Includes) == 0 && len(expected.Items.Excludes) == 0)
}

func verifyCompletionsItemDefaults(t *testing.T, actual *lsproto.CompletionItemDefaults, expected *CompletionsExpectedItemDefaults, prefix string) {
	if actual == nil {
		if expected == nil {
			return
		}
		t.Fatalf(prefix+"Expected non-nil completion item defaults but got nil: %s", cmp.Diff(actual, nil))
	}
	if expected == nil {
		t.Fatalf(prefix+"Expected nil completion item defaults but got non-nil: %s", cmp.Diff(actual, nil))
	}
	assertDeepEqual(t, actual.CommitCharacters, expected.CommitCharacters, prefix+"CommitCharacters mismatch:")
	switch editRange := expected.EditRange.(type) {
	case *EditRange:
		if actual.EditRange == nil {
			t.Fatal(prefix + "Expected non-nil EditRange but got nil")
		}
		expectedInsert := editRange.Insert.LSRange
		expectedReplace := editRange.Replace.LSRange
		assertDeepEqual(
			t,
			actual.EditRange,
			&lsproto.RangeOrEditRangeWithInsertReplace{
				EditRangeWithInsertReplace: &lsproto.EditRangeWithInsertReplace{
					Insert:  expectedInsert,
					Replace: expectedReplace,
				},
			},
			prefix+"EditRange mismatch:")
	case nil:
		if actual.EditRange != nil {
			t.Fatalf(prefix+"Expected nil EditRange but got non-nil: %s", cmp.Diff(actual.EditRange, nil))
		}
	case Ignored:
	default:
		t.Fatalf(prefix+"Expected EditRange to be *EditRange or Ignored, got %T", editRange)
	}
}

func verifyCompletionsItems(t *testing.T, prefix string, actual []*lsproto.CompletionItem, expected *CompletionsExpectedItems) {
	if expected.Exact != nil {
		if expected.Includes != nil {
			t.Fatal(prefix + "Expected exact completion list but also specified 'includes'.")
		}
		if expected.Excludes != nil {
			t.Fatal(prefix + "Expected exact completion list but also specified 'excludes'.")
		}
		if len(actual) != len(expected.Exact) {
			t.Fatalf(prefix+"Expected %d exact completion items but got %d: %s", len(expected.Exact), len(actual), cmp.Diff(actual, expected.Exact))
		}
		if len(actual) > 0 {
			verifyCompletionsAreExactly(t, prefix, actual, expected.Exact)
		}
		return
	}
	nameToActualItem := make(map[string]*lsproto.CompletionItem)
	for _, item := range actual {
		nameToActualItem[item.Label] = item
	}
	if expected.Includes != nil {
		for _, item := range expected.Includes {
			switch item := item.(type) {
			case string:
				_, ok := nameToActualItem[item]
				if !ok {
					t.Fatalf("%sLabel '%s' not found in actual items. Actual items: %s", prefix, item, cmp.Diff(actual, nil))
				}
			case *lsproto.CompletionItem:
				actualItem, ok := nameToActualItem[item.Label]
				if !ok {
					t.Fatalf("%sLabel '%s' not found in actual items. Actual items: %s", prefix, item.Label, cmp.Diff(actual, nil))
				}
				verifyCompletionItem(t, prefix+"Includes completion item mismatch for label "+item.Label, actualItem, item)
			default:
				t.Fatalf("%sExpected completion item to be a string or *lsproto.CompletionItem, got %T", prefix, item)
			}
		}
	}
	for _, exclude := range expected.Excludes {
		if _, ok := nameToActualItem[exclude]; ok {
			t.Fatalf("%sLabel '%s' should not be in actual items but was found. Actual items: %s", prefix, exclude, cmp.Diff(actual, nil))
		}
	}
}

func verifyCompletionsAreExactly(t *testing.T, prefix string, actual []*lsproto.CompletionItem, expected []CompletionsExpectedItem) {
	// Verify labels first
	assertDeepEqual(t, core.Map(actual, func(item *lsproto.CompletionItem) string {
		return item.Label
	}), core.Map(expected, func(item CompletionsExpectedItem) string {
		return getExpectedLabel(t, item)
	}), prefix+"Labels mismatch")
	for i, actualItem := range actual {
		switch expectedItem := expected[i].(type) {
		case string:
			continue // already checked labels
		case *lsproto.CompletionItem:
			verifyCompletionItem(t, prefix+"Completion item mismatch for label "+actualItem.Label, actualItem, expectedItem)
		}
	}
}

func verifyCompletionItem(t *testing.T, prefix string, actual *lsproto.CompletionItem, expected *lsproto.CompletionItem) {
	ignoreKind := cmp.FilterPath(
		func(p cmp.Path) bool {
			switch p.Last().String() {
			case ".Kind", ".SortText":
				return true
			default:
				return false
			}
		},
		cmp.Ignore(),
	)
	assertDeepEqual(t, actual, expected, prefix, ignoreKind)
	if expected.Kind != nil {
		assertDeepEqual(t, actual.Kind, expected.Kind, prefix+" Kind mismatch")
	}
	assertDeepEqual(t, actual.SortText, core.OrElse(expected.SortText, ptrTo(string(ls.SortTextLocationPriority))), prefix+" SortText mismatch")
}

func getExpectedLabel(t *testing.T, item CompletionsExpectedItem) string {
	switch item := item.(type) {
	case string:
		return item
	case *lsproto.CompletionItem:
		return item.Label
	default:
		t.Fatalf("Expected completion item to be a string or *lsproto.CompletionItem, got %T", item)
		return ""
	}
}

func assertDeepEqual(t *testing.T, actual any, expected any, prefix string, opts ...cmp.Option) {
	t.Helper()

	diff := cmp.Diff(actual, expected, opts...)
	if diff != "" {
		t.Fatalf("%s:\n%s", prefix, diff)
	}
}

func ptrTo[T any](v T) *T {
	return &v
}
