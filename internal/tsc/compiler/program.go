package compiler

import (
	"context"
	"maps"
	"slices"
	"sync"

	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/binder"
	"github.com/yasufadhili/jawt/internal/tsc/checker"
	"github.com/yasufadhili/jawt/internal/tsc/collections"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/diagnostics"
	"github.com/yasufadhili/jawt/internal/tsc/module"
	"github.com/yasufadhili/jawt/internal/tsc/modulespecifiers"
	"github.com/yasufadhili/jawt/internal/tsc/outputpaths"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/scanner"
	"github.com/yasufadhili/jawt/internal/tsc/sourcemap"
	"github.com/yasufadhili/jawt/internal/tsc/tsoptions"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
)

type ProgramOptions struct {
	Host                        CompilerHost
	Config                      *tsoptions.ParsedCommandLine
	UseSourceOfProjectReference bool
	SingleThreaded              core.Tristate
	CreateCheckerPool           func(*Program) CheckerPool
	TypingsLocation             string
	ProjectName                 string
	JSDocParsingMode            ast.JSDocParsingMode
}

func (p *ProgramOptions) canUseProjectReferenceSource() bool {
	return p.UseSourceOfProjectReference && !p.Config.CompilerOptions().DisableSourceOfProjectReferenceRedirect.IsTrue()
}

type Program struct {
	opts        ProgramOptions
	nodeModules map[string]*ast.SourceFile
	checkerPool CheckerPool

	comparePathsOptions tspath.ComparePathsOptions

	processedFiles

	usesUriStyleNodeCoreModules core.Tristate

	commonSourceDirectory     string
	commonSourceDirectoryOnce sync.Once

	declarationDiagnosticCache collections.SyncMap[*ast.SourceFile, []*ast.Diagnostic]
}

// FileExists implements checker.Program.
func (p *Program) FileExists(path string) bool {
	return p.Host().FS().FileExists(path)
}

// GetCurrentDirectory implements checker.Program.
func (p *Program) GetCurrentDirectory() string {
	return p.Host().GetCurrentDirectory()
}

// GetGlobalTypingsCacheLocation implements checker.Program.
func (p *Program) GetGlobalTypingsCacheLocation() string {
	return "" // !!! see src/tsserver/nodeServer.ts for strada's node-specific implementation
}

// GetNearestAncestorDirectoryWithPackageJson implements checker.Program.
func (p *Program) GetNearestAncestorDirectoryWithPackageJson(dirname string) string {
	scoped := p.resolver.GetPackageScopeForPath(dirname)
	if scoped != nil && scoped.Exists() {
		return scoped.PackageDirectory
	}
	return ""
}

// GetPackageJsonInfo implements checker.Program.
func (p *Program) GetPackageJsonInfo(pkgJsonPath string) modulespecifiers.PackageJsonInfo {
	scoped := p.resolver.GetPackageScopeForPath(pkgJsonPath)
	if scoped != nil && scoped.Exists() && scoped.PackageDirectory == tspath.GetDirectoryPath(pkgJsonPath) {
		return scoped
	}
	return nil
}

// GetRedirectTargets implements checker.Program.
func (p *Program) GetRedirectTargets(path tspath.Path) []string {
	return nil // !!! TODO: project references support
}

// GetOutputAndProjectReference implements checker.Program.
func (p *Program) GetOutputAndProjectReference(path tspath.Path) *tsoptions.OutputDtsAndProjectReference {
	return p.projectReferenceFileMapper.getOutputAndProjectReference(path)
}

// IsSourceFromProjectReference implements checker.Program.
func (p *Program) IsSourceFromProjectReference(path tspath.Path) bool {
	return p.projectReferenceFileMapper.isSourceFromProjectReference(path)
}

func (p *Program) GetSourceAndProjectReference(path tspath.Path) *tsoptions.SourceAndProjectReference {
	return p.projectReferenceFileMapper.getSourceAndProjectReference(path)
}

func (p *Program) GetResolvedProjectReferenceFor(path tspath.Path) (*tsoptions.ParsedCommandLine, bool) {
	return p.projectReferenceFileMapper.getResolvedReferenceFor(path)
}

func (p *Program) GetRedirectForResolution(file ast.HasFileName) *tsoptions.ParsedCommandLine {
	return p.projectReferenceFileMapper.getRedirectForResolution(file)
}

func (p *Program) ForEachResolvedProjectReference(
	fn func(path tspath.Path, config *tsoptions.ParsedCommandLine),
) {
	p.projectReferenceFileMapper.forEachResolvedProjectReference(fn)
}

// UseCaseSensitiveFileNames implements checker.Program.
func (p *Program) UseCaseSensitiveFileNames() bool {
	return p.Host().FS().UseCaseSensitiveFileNames()
}

var _ checker.Program = (*Program)(nil)

/** This should have similar behavior to 'processSourceFile' without diagnostics or mutation. */
func (p *Program) GetSourceFileFromReference(origin *ast.SourceFile, ref *ast.FileReference) *ast.SourceFile {
	// TODO: The module loader in corsa is fairly different than strada, it should probably be able to expose this functionality at some point,
	// rather than redoing the logic approximately here, since most of the related logic now lives in module.Resolver
	// Still, without the failed lookup reporting that only the loader does, this isn't terribly complicated

	fileName := tspath.ResolvePath(tspath.GetDirectoryPath(origin.FileName()), ref.FileName)
	supportedExtensionsBase := tsoptions.GetSupportedExtensions(p.Options(), nil /*extraFileExtensions*/)
	supportedExtensions := tsoptions.GetSupportedExtensionsWithJsonIfResolveJsonModule(p.Options(), supportedExtensionsBase)
	allowNonTsExtensions := p.Options().AllowNonTsExtensions.IsTrue()
	if tspath.HasExtension(fileName) {
		if !allowNonTsExtensions {
			canonicalFileName := tspath.GetCanonicalFileName(fileName, p.UseCaseSensitiveFileNames())
			supported := false
			for _, group := range supportedExtensions {
				if tspath.FileExtensionIsOneOf(canonicalFileName, group) {
					supported = true
					break
				}
			}
			if !supported {
				return nil // unsupported extensions are forced to fail
			}
		}

		return p.GetSourceFile(fileName)
	}
	if allowNonTsExtensions {
		extensionless := p.GetSourceFile(fileName)
		if extensionless != nil {
			return extensionless
		}
	}

	// Only try adding extensions from the first supported group (which should be .ts/.tsx/.d.ts)
	for _, ext := range supportedExtensions[0] {
		result := p.GetSourceFile(fileName + ext)
		if result != nil {
			return result
		}
	}
	return nil
}

func NewProgram(opts ProgramOptions) *Program {
	p := &Program{opts: opts}
	p.initCheckerPool()
	p.processedFiles = processAllProgramFiles(p.opts, p.singleThreaded())
	return p
}

// Return an updated program for which it is known that only the file with the given path has changed.
// In addition to a new program, return a boolean indicating whether the data of the old program was reused.
func (p *Program) UpdateProgram(changedFilePath tspath.Path) (*Program, bool) {
	oldFile := p.filesByPath[changedFilePath]
	newFile := p.Host().GetSourceFile(oldFile.ParseOptions())
	if !canReplaceFileInProgram(oldFile, newFile) {
		return NewProgram(p.opts), false
	}
	result := &Program{
		opts:                        p.opts,
		nodeModules:                 p.nodeModules,
		comparePathsOptions:         p.comparePathsOptions,
		processedFiles:              p.processedFiles,
		usesUriStyleNodeCoreModules: p.usesUriStyleNodeCoreModules,
	}
	result.initCheckerPool()
	index := core.FindIndex(result.files, func(file *ast.SourceFile) bool { return file.Path() == newFile.Path() })
	result.files = slices.Clone(result.files)
	result.files[index] = newFile
	result.filesByPath = maps.Clone(result.filesByPath)
	result.filesByPath[newFile.Path()] = newFile
	return result, true
}

func (p *Program) initCheckerPool() {
	if p.opts.CreateCheckerPool != nil {
		p.checkerPool = p.opts.CreateCheckerPool(p)
	} else {
		p.checkerPool = newCheckerPool(core.IfElse(p.singleThreaded(), 1, 4), p)
	}
}

func canReplaceFileInProgram(file1 *ast.SourceFile, file2 *ast.SourceFile) bool {
	return file2 != nil &&
		file1.ParseOptions() == file2.ParseOptions() &&
		file1.UsesUriStyleNodeCoreModules == file2.UsesUriStyleNodeCoreModules &&
		slices.EqualFunc(file1.Imports(), file2.Imports(), equalModuleSpecifiers) &&
		slices.EqualFunc(file1.ModuleAugmentations, file2.ModuleAugmentations, equalModuleAugmentationNames) &&
		slices.Equal(file1.AmbientModuleNames, file2.AmbientModuleNames) &&
		slices.EqualFunc(file1.ReferencedFiles, file2.ReferencedFiles, equalFileReferences) &&
		slices.EqualFunc(file1.TypeReferenceDirectives, file2.TypeReferenceDirectives, equalFileReferences) &&
		slices.EqualFunc(file1.LibReferenceDirectives, file2.LibReferenceDirectives, equalFileReferences) &&
		equalCheckJSDirectives(file1.CheckJsDirective, file2.CheckJsDirective)
}

func equalModuleSpecifiers(n1 *ast.Node, n2 *ast.Node) bool {
	return n1.Kind == n2.Kind && (!ast.IsStringLiteral(n1) || n1.Text() == n2.Text())
}

func equalModuleAugmentationNames(n1 *ast.Node, n2 *ast.Node) bool {
	return n1.Kind == n2.Kind && n1.Text() == n2.Text()
}

func equalFileReferences(f1 *ast.FileReference, f2 *ast.FileReference) bool {
	return f1.FileName == f2.FileName && f1.ResolutionMode == f2.ResolutionMode && f1.Preserve == f2.Preserve
}

func equalCheckJSDirectives(d1 *ast.CheckJsDirective, d2 *ast.CheckJsDirective) bool {
	return d1 == nil && d2 == nil || d1 != nil && d2 != nil && d1.Enabled == d2.Enabled
}

func (p *Program) SourceFiles() []*ast.SourceFile { return p.files }
func (p *Program) Options() *core.CompilerOptions { return p.opts.Config.CompilerOptions() }
func (p *Program) Host() CompilerHost             { return p.opts.Host }
func (p *Program) GetConfigFileParsingDiagnostics() []*ast.Diagnostic {
	return slices.Clip(p.opts.Config.GetConfigFileParsingDiagnostics())
}

func (p *Program) singleThreaded() bool {
	return p.opts.SingleThreaded.DefaultIfUnknown(p.Options().SingleThreaded).IsTrue()
}

func (p *Program) BindSourceFiles() {
	wg := core.NewWorkGroup(p.singleThreaded())
	for _, file := range p.files {
		if !file.IsBound() {
			wg.Queue(func() {
				binder.BindSourceFile(file)
			})
		}
	}
	wg.RunAndWait()
}

func (p *Program) CheckSourceFiles(ctx context.Context) {
	wg := core.NewWorkGroup(p.singleThreaded())
	checkers, done := p.checkerPool.GetAllCheckers(ctx)
	defer done()
	for _, checker := range checkers {
		wg.Queue(func() {
			for file := range p.checkerPool.Files(checker) {
				checker.CheckSourceFile(ctx, file)
			}
		})
	}
	wg.RunAndWait()
}

// Return the type checker associated with the program.
func (p *Program) GetTypeChecker(ctx context.Context) (*checker.Checker, func()) {
	return p.checkerPool.GetChecker(ctx)
}

func (p *Program) GetTypeCheckers(ctx context.Context) ([]*checker.Checker, func()) {
	return p.checkerPool.GetAllCheckers(ctx)
}

// Return a checker for the given file. We may have multiple checkers in concurrent scenarios and this
// method returns the checker that was tasked with checking the file. Note that it isn't possible to mix
// types obtained from different checkers, so only non-type data (such as diagnostics or string
// representations of types) should be obtained from checkers returned by this method.
func (p *Program) GetTypeCheckerForFile(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func()) {
	return p.checkerPool.GetCheckerForFile(ctx, file)
}

func (p *Program) GetResolvedModule(file ast.HasFileName, moduleReference string, mode core.ResolutionMode) *module.ResolvedModule {
	if resolutions, ok := p.resolvedModules[file.Path()]; ok {
		if resolved, ok := resolutions[module.ModeAwareCacheKey{Name: moduleReference, Mode: mode}]; ok {
			return resolved
		}
	}
	return nil
}

func (p *Program) GetResolvedModuleFromModuleSpecifier(file ast.HasFileName, moduleSpecifier *ast.StringLiteralLike) *module.ResolvedModule {
	if !ast.IsStringLiteralLike(moduleSpecifier) {
		panic("moduleSpecifier must be a StringLiteralLike")
	}
	mode := p.GetModeForUsageLocation(file, moduleSpecifier)
	return p.GetResolvedModule(file, moduleSpecifier.Text(), mode)
}

func (p *Program) GetResolvedModules() map[tspath.Path]module.ModeAwareCache[*module.ResolvedModule] {
	return p.resolvedModules
}

func (p *Program) GetSyntacticDiagnostics(ctx context.Context, sourceFile *ast.SourceFile) []*ast.Diagnostic {
	return p.getDiagnosticsHelper(ctx, sourceFile, false /*ensureBound*/, false /*ensureChecked*/, p.getSyntacticDiagnosticsForFile)
}

func (p *Program) GetBindDiagnostics(ctx context.Context, sourceFile *ast.SourceFile) []*ast.Diagnostic {
	return p.getDiagnosticsHelper(ctx, sourceFile, true /*ensureBound*/, false /*ensureChecked*/, p.getBindDiagnosticsForFile)
}

func (p *Program) GetSemanticDiagnostics(ctx context.Context, sourceFile *ast.SourceFile) []*ast.Diagnostic {
	return p.getDiagnosticsHelper(ctx, sourceFile, true /*ensureBound*/, true /*ensureChecked*/, p.getSemanticDiagnosticsForFile)
}

func (p *Program) GetSuggestionDiagnostics(ctx context.Context, sourceFile *ast.SourceFile) []*ast.Diagnostic {
	return p.getDiagnosticsHelper(ctx, sourceFile, true /*ensureBound*/, true /*ensureChecked*/, p.getSuggestionDiagnosticsForFile)
}

func (p *Program) GetGlobalDiagnostics(ctx context.Context) []*ast.Diagnostic {
	var globalDiagnostics []*ast.Diagnostic
	checkers, done := p.checkerPool.GetAllCheckers(ctx)
	defer done()
	for _, checker := range checkers {
		globalDiagnostics = append(globalDiagnostics, checker.GetGlobalDiagnostics()...)
	}

	return SortAndDeduplicateDiagnostics(globalDiagnostics)
}

func (p *Program) GetDeclarationDiagnostics(ctx context.Context, sourceFile *ast.SourceFile) []*ast.Diagnostic {
	return p.getDiagnosticsHelper(ctx, sourceFile, true /*ensureBound*/, true /*ensureChecked*/, p.getDeclarationDiagnosticsForFile)
}

func (p *Program) GetOptionsDiagnostics(ctx context.Context) []*ast.Diagnostic {
	return SortAndDeduplicateDiagnostics(append(p.GetGlobalDiagnostics(ctx), p.getOptionsDiagnosticsOfConfigFile()...))
}

func (p *Program) getOptionsDiagnosticsOfConfigFile() []*ast.Diagnostic {
	// todo update p.configParsingDiagnostics when updateAndGetProgramDiagnostics is implemented
	if p.Options() == nil || p.Options().ConfigFilePath == "" {
		return nil
	}
	return p.GetConfigFileParsingDiagnostics() // TODO: actually call getDiagnosticsHelper on config path
}

func (p *Program) getSyntacticDiagnosticsForFile(ctx context.Context, sourceFile *ast.SourceFile) []*ast.Diagnostic {
	return sourceFile.Diagnostics()
}

func (p *Program) getBindDiagnosticsForFile(ctx context.Context, sourceFile *ast.SourceFile) []*ast.Diagnostic {
	// TODO: restore this; tsgo's main depends on this function binding all files for timing.
	// if checker.SkipTypeChecking(sourceFile, p.compilerOptions) {
	// 	return nil
	// }

	return sourceFile.BindDiagnostics()
}

func (p *Program) getSemanticDiagnosticsForFile(ctx context.Context, sourceFile *ast.SourceFile) []*ast.Diagnostic {
	compilerOptions := p.Options()
	if checker.SkipTypeChecking(sourceFile, compilerOptions, p) {
		return nil
	}

	var fileChecker *checker.Checker
	var done func()
	if sourceFile != nil {
		fileChecker, done = p.checkerPool.GetCheckerForFile(ctx, sourceFile)
		defer done()
	}
	diags := slices.Clip(sourceFile.BindDiagnostics())
	checkers, closeCheckers := p.checkerPool.GetAllCheckers(ctx)
	defer closeCheckers()

	// Ask for diags from all checkers; checking one file may add diagnostics to other files.
	// These are deduplicated later.
	for _, checker := range checkers {
		if sourceFile == nil || checker == fileChecker {
			diags = append(diags, checker.GetDiagnostics(ctx, sourceFile)...)
		} else {
			diags = append(diags, checker.GetDiagnosticsWithoutCheck(sourceFile)...)
		}
	}
	if ctx.Err() != nil {
		return nil
	}

	// !!! This should be rewritten to work like getBindAndCheckDiagnosticsForFileNoCache.

	isPlainJS := ast.IsPlainJSFile(sourceFile, compilerOptions.CheckJs)
	if isPlainJS {
		return core.Filter(diags, func(d *ast.Diagnostic) bool {
			return plainJSErrors.Has(d.Code())
		})
	}

	if len(sourceFile.CommentDirectives) == 0 {
		return diags
	}
	// Build map of directives by line number
	directivesByLine := make(map[int]ast.CommentDirective)
	for _, directive := range sourceFile.CommentDirectives {
		line, _ := scanner.GetLineAndCharacterOfPosition(sourceFile, directive.Loc.Pos())
		directivesByLine[line] = directive
	}
	lineStarts := scanner.GetLineStarts(sourceFile)
	filtered := make([]*ast.Diagnostic, 0, len(diags))
	for _, diagnostic := range diags {
		ignoreDiagnostic := false
		for line := scanner.ComputeLineOfPosition(lineStarts, diagnostic.Pos()) - 1; line >= 0; line-- {
			// If line contains a @ts-ignore or @ts-expect-error directive, ignore this diagnostic and change
			// the directive kind to @ts-ignore to indicate it was used.
			if directive, ok := directivesByLine[line]; ok {
				ignoreDiagnostic = true
				directive.Kind = ast.CommentDirectiveKindIgnore
				directivesByLine[line] = directive
				break
			}
			// Stop searching backwards when we encounter a line that isn't blank or a comment.
			if !isCommentOrBlankLine(sourceFile.Text(), int(lineStarts[line])) {
				break
			}
		}
		if !ignoreDiagnostic {
			filtered = append(filtered, diagnostic)
		}
	}
	for _, directive := range directivesByLine {
		// Above we changed all used directive kinds to @ts-ignore, so any @ts-expect-error directives that
		// remain are unused and thus errors.
		if directive.Kind == ast.CommentDirectiveKindExpectError {
			filtered = append(filtered, ast.NewDiagnostic(sourceFile, directive.Loc, diagnostics.Unused_ts_expect_error_directive))
		}
	}
	return filtered
}

func (p *Program) getDeclarationDiagnosticsForFile(ctx context.Context, sourceFile *ast.SourceFile) []*ast.Diagnostic {
	if sourceFile.IsDeclarationFile {
		return []*ast.Diagnostic{}
	}

	if cached, ok := p.declarationDiagnosticCache.Load(sourceFile); ok {
		return cached
	}

	host, done := newEmitHost(ctx, p, sourceFile)
	defer done()
	diagnostics := getDeclarationDiagnostics(host, sourceFile)
	diagnostics, _ = p.declarationDiagnosticCache.LoadOrStore(sourceFile, diagnostics)
	return diagnostics
}

func (p *Program) getSuggestionDiagnosticsForFile(ctx context.Context, sourceFile *ast.SourceFile) []*ast.Diagnostic {
	if checker.SkipTypeChecking(sourceFile, p.Options(), p) {
		return nil
	}

	var fileChecker *checker.Checker
	var done func()
	if sourceFile != nil {
		fileChecker, done = p.checkerPool.GetCheckerForFile(ctx, sourceFile)
		defer done()
	}

	diags := slices.Clip(sourceFile.BindSuggestionDiagnostics)

	checkers, closeCheckers := p.checkerPool.GetAllCheckers(ctx)
	defer closeCheckers()

	// Ask for diags from all checkers; checking one file may add diagnostics to other files.
	// These are deduplicated later.
	for _, checker := range checkers {
		if sourceFile == nil || checker == fileChecker {
			diags = append(diags, checker.GetSuggestionDiagnostics(ctx, sourceFile)...)
		} else {
			// !!! is there any case where suggestion diagnostics are produced in other checkers?
		}
	}
	if ctx.Err() != nil {
		return nil
	}

	return diags
}

func isCommentOrBlankLine(text string, pos int) bool {
	for pos < len(text) && (text[pos] == ' ' || text[pos] == '\t') {
		pos++
	}
	return pos == len(text) ||
		pos < len(text) && (text[pos] == '\r' || text[pos] == '\n') ||
		pos+1 < len(text) && text[pos] == '/' && text[pos+1] == '/'
}

func SortAndDeduplicateDiagnostics(diagnostics []*ast.Diagnostic) []*ast.Diagnostic {
	diagnostics = slices.Clone(diagnostics)
	slices.SortFunc(diagnostics, ast.CompareDiagnostics)
	return compactAndMergeRelatedInfos(diagnostics)
}

// Remove duplicate diagnostics and, for sequences of diagnostics that differ only by related information,
// create a single diagnostic with sorted and deduplicated related information.
func compactAndMergeRelatedInfos(diagnostics []*ast.Diagnostic) []*ast.Diagnostic {
	if len(diagnostics) < 2 {
		return diagnostics
	}
	i := 0
	j := 0
	for i < len(diagnostics) {
		d := diagnostics[i]
		n := 1
		for i+n < len(diagnostics) && ast.EqualDiagnosticsNoRelatedInfo(d, diagnostics[i+n]) {
			n++
		}
		if n > 1 {
			var relatedInfos []*ast.Diagnostic
			for k := range n {
				relatedInfos = append(relatedInfos, diagnostics[i+k].RelatedInformation()...)
			}
			if relatedInfos != nil {
				slices.SortFunc(relatedInfos, ast.CompareDiagnostics)
				relatedInfos = slices.CompactFunc(relatedInfos, ast.EqualDiagnostics)
				d = d.Clone().SetRelatedInfo(relatedInfos)
			}
		}
		diagnostics[j] = d
		i += n
		j++
	}
	clear(diagnostics[j:])
	return diagnostics[:j]
}

func (p *Program) getDiagnosticsHelper(ctx context.Context, sourceFile *ast.SourceFile, ensureBound bool, ensureChecked bool, getDiagnostics func(context.Context, *ast.SourceFile) []*ast.Diagnostic) []*ast.Diagnostic {
	if sourceFile != nil {
		if ensureBound {
			binder.BindSourceFile(sourceFile)
		}
		return SortAndDeduplicateDiagnostics(getDiagnostics(ctx, sourceFile))
	}
	if ensureBound {
		p.BindSourceFiles()
	}
	if ensureChecked {
		p.CheckSourceFiles(ctx)
		if ctx.Err() != nil {
			return nil
		}
	}
	var result []*ast.Diagnostic
	for _, file := range p.files {
		result = append(result, getDiagnostics(ctx, file)...)
	}
	return SortAndDeduplicateDiagnostics(result)
}

func (p *Program) LineCount() int {
	var count int
	for _, file := range p.files {
		count += len(file.LineMap())
	}
	return count
}

func (p *Program) IdentifierCount() int {
	var count int
	for _, file := range p.files {
		count += file.IdentifierCount
	}
	return count
}

func (p *Program) SymbolCount() int {
	var count int
	for _, file := range p.files {
		count += file.SymbolCount
	}
	checkers, done := p.checkerPool.GetAllCheckers(context.Background())
	defer done()
	for _, checker := range checkers {
		count += int(checker.SymbolCount)
	}
	return count
}

func (p *Program) TypeCount() int {
	var count int
	checkers, done := p.checkerPool.GetAllCheckers(context.Background())
	defer done()
	for _, checker := range checkers {
		count += int(checker.TypeCount)
	}
	return count
}

func (p *Program) InstantiationCount() int {
	var count int
	checkers, done := p.checkerPool.GetAllCheckers(context.Background())
	defer done()
	for _, checker := range checkers {
		count += int(checker.TotalInstantiationCount)
	}
	return count
}

func (p *Program) GetSourceFileMetaData(path tspath.Path) ast.SourceFileMetaData {
	return p.sourceFileMetaDatas[path]
}

func (p *Program) GetEmitModuleFormatOfFile(sourceFile ast.HasFileName) core.ModuleKind {
	return ast.GetEmitModuleFormatOfFileWorker(sourceFile.FileName(), p.projectReferenceFileMapper.getCompilerOptionsForFile(sourceFile), p.GetSourceFileMetaData(sourceFile.Path()))
}

func (p *Program) GetEmitSyntaxForUsageLocation(sourceFile ast.HasFileName, location *ast.StringLiteralLike) core.ResolutionMode {
	return getEmitSyntaxForUsageLocationWorker(sourceFile.FileName(), p.sourceFileMetaDatas[sourceFile.Path()], location, p.projectReferenceFileMapper.getCompilerOptionsForFile(sourceFile))
}

func (p *Program) GetImpliedNodeFormatForEmit(sourceFile ast.HasFileName) core.ResolutionMode {
	return ast.GetImpliedNodeFormatForEmitWorker(sourceFile.FileName(), p.projectReferenceFileMapper.getCompilerOptionsForFile(sourceFile).GetEmitModuleKind(), p.GetSourceFileMetaData(sourceFile.Path()))
}

func (p *Program) GetModeForUsageLocation(sourceFile ast.HasFileName, location *ast.StringLiteralLike) core.ResolutionMode {
	return getModeForUsageLocation(sourceFile.FileName(), p.sourceFileMetaDatas[sourceFile.Path()], location, p.projectReferenceFileMapper.getCompilerOptionsForFile(sourceFile))
}

func (p *Program) GetDefaultResolutionModeForFile(sourceFile ast.HasFileName) core.ResolutionMode {
	return getDefaultResolutionModeForFile(sourceFile.FileName(), p.sourceFileMetaDatas[sourceFile.Path()], p.projectReferenceFileMapper.getCompilerOptionsForFile(sourceFile))
}

func (p *Program) IsSourceFileDefaultLibrary(path tspath.Path) bool {
	return p.libFiles.Has(path)
}

func (p *Program) CommonSourceDirectory() string {
	p.commonSourceDirectoryOnce.Do(func() {
		p.commonSourceDirectory = outputpaths.GetCommonSourceDirectory(
			p.Options(),
			func() []string {
				var files []string
				for _, file := range p.files {
					if sourceFileMayBeEmitted(file, p, false /*forceDtsEmit*/) {
						files = append(files, file.FileName())
					}
				}
				return files
			},
			p.GetCurrentDirectory(),
			p.UseCaseSensitiveFileNames(),
		)
	})
	return p.commonSourceDirectory
}

type EmitOptions struct {
	TargetSourceFile *ast.SourceFile // Single file to emit. If `nil`, emits all files
	forceDtsEmit     bool
}

type EmitResult struct {
	EmitSkipped  bool
	Diagnostics  []*ast.Diagnostic      // Contains declaration emit diagnostics
	EmittedFiles []string               // Array of files the compiler wrote to disk
	SourceMaps   []*SourceMapEmitResult // Array of sourceMapData if compiler emitted sourcemaps
}

type SourceMapEmitResult struct {
	InputSourceFileNames []string // Input source file (which one can use on program to get the file), 1:1 mapping with the sourceMap.sources list
	SourceMap            *sourcemap.RawSourceMap
	GeneratedFile        string
}

func (p *Program) Emit(options EmitOptions) *EmitResult {
	// !!! performance measurement
	p.BindSourceFiles()

	writerPool := &sync.Pool{
		New: func() any {
			return printer.NewTextWriter(p.Options().NewLine.GetNewLineCharacter())
		},
	}
	wg := core.NewWorkGroup(p.singleThreaded())
	var emitters []*emitter
	sourceFiles := getSourceFilesToEmit(p, options.TargetSourceFile, options.forceDtsEmit)

	for _, sourceFile := range sourceFiles {
		emitter := &emitter{
			emittedFilesList:  nil,
			sourceMapDataList: nil,
			writer:            nil,
			sourceFile:        sourceFile,
		}
		emitters = append(emitters, emitter)
		wg.Queue(func() {
			host, done := newEmitHost(context.TODO(), p, sourceFile)
			defer done()
			emitter.host = host

			// take an unused writer
			writer := writerPool.Get().(printer.EmitTextWriter)
			writer.Clear()

			// attach writer and perform emit
			emitter.writer = writer
			emitter.paths = outputpaths.GetOutputPathsFor(sourceFile, host.Options(), host, options.forceDtsEmit)
			emitter.emit()
			emitter.writer = nil

			// put the writer back in the pool
			writerPool.Put(writer)
		})
	}

	// wait for emit to complete
	wg.RunAndWait()

	// collect results from emit, preserving input order
	result := &EmitResult{}
	for _, emitter := range emitters {
		if emitter.emitSkipped {
			result.EmitSkipped = true
		}
		result.Diagnostics = append(result.Diagnostics, emitter.emitterDiagnostics.GetDiagnostics()...)
		if emitter.emittedFilesList != nil {
			result.EmittedFiles = append(result.EmittedFiles, emitter.emittedFilesList...)
		}
		if emitter.sourceMapDataList != nil {
			result.SourceMaps = append(result.SourceMaps, emitter.sourceMapDataList...)
		}
	}
	return result
}

func (p *Program) GetSourceFile(filename string) *ast.SourceFile {
	path := tspath.ToPath(filename, p.GetCurrentDirectory(), p.UseCaseSensitiveFileNames())
	return p.GetSourceFileByPath(path)
}

func (p *Program) GetSourceFileForResolvedModule(fileName string) *ast.SourceFile {
	file := p.GetSourceFile(fileName)
	if file == nil {
		filename := p.projectReferenceFileMapper.getParseFileRedirect(ast.NewHasFileName(fileName, tspath.ToPath(fileName, p.GetCurrentDirectory(), p.UseCaseSensitiveFileNames())))
		if filename != "" {
			return p.GetSourceFile(filename)
		}
	}
	return file
}

func (p *Program) GetSourceFileByPath(path tspath.Path) *ast.SourceFile {
	return p.filesByPath[path]
}

func (p *Program) GetSourceFiles() []*ast.SourceFile {
	return p.files
}

func (p *Program) GetLibFileFromReference(ref *ast.FileReference) *ast.SourceFile {
	path, ok := tsoptions.GetLibFileName(ref.FileName)
	if !ok {
		return nil
	}
	if sourceFile, ok := p.filesByPath[tspath.Path(path)]; ok {
		return sourceFile
	}
	return nil
}

func (p *Program) GetResolvedTypeReferenceDirectiveFromTypeReferenceDirective(typeRef *ast.FileReference, sourceFile *ast.SourceFile) *module.ResolvedTypeReferenceDirective {
	if resolutions, ok := p.typeResolutionsInFile[sourceFile.Path()]; ok {
		if resolved, ok := resolutions[module.ModeAwareCacheKey{Name: typeRef.FileName, Mode: p.getModeForTypeReferenceDirectiveInFile(typeRef, sourceFile)}]; ok {
			return resolved
		}
	}
	return nil
}

func (p *Program) GetResolvedTypeReferenceDirectives() map[tspath.Path]module.ModeAwareCache[*module.ResolvedTypeReferenceDirective] {
	return p.typeResolutionsInFile
}

func (p *Program) getModeForTypeReferenceDirectiveInFile(ref *ast.FileReference, sourceFile *ast.SourceFile) core.ResolutionMode {
	if ref.ResolutionMode != core.ResolutionModeNone {
		return ref.ResolutionMode
	}
	return p.GetDefaultResolutionModeForFile(sourceFile)
}

func (p *Program) IsSourceFileFromExternalLibrary(file *ast.SourceFile) bool {
	return p.sourceFilesFoundSearchingNodeModules.Has(file.Path())
}

type FileIncludeKind int

const (
	FileIncludeKindRootFile FileIncludeKind = iota
	FileIncludeKindSourceFromProjectReference
	FileIncludeKindOutputFromProjectReference
	FileIncludeKindImport
	FileIncludeKindReferenceFile
	FileIncludeKindTypeReferenceDirective
	FileIncludeKindLibFile
	FileIncludeKindLibReferenceDirective
	FileIncludeKindAutomaticTypeDirectiveFile
)

type FileIncludeReason struct {
	Kind  FileIncludeKind
	Index int
}

// UnsupportedExtensions returns a list of all present "unsupported" extensions,
// e.g. extensions that are not yet supported by the port.
func (p *Program) UnsupportedExtensions() []string {
	return p.unsupportedExtensions
}

func (p *Program) GetJSXRuntimeImportSpecifier(path tspath.Path) (moduleReference string, specifier *ast.Node) {
	if result := p.jsxRuntimeImportSpecifiers[path]; result != nil {
		return result.moduleReference, result.specifier
	}
	return "", nil
}

func (p *Program) GetImportHelpersImportSpecifier(path tspath.Path) *ast.Node {
	return p.importHelpersImportSpecifiers[path]
}

func (p *Program) SourceFileMayBeEmitted(sourceFile *ast.SourceFile, forceDtsEmit bool) bool {
	return sourceFileMayBeEmitted(sourceFile, &emitHost{program: p}, forceDtsEmit)
}

var plainJSErrors = collections.NewSetFromItems(
	// binder errors
	diagnostics.Cannot_redeclare_block_scoped_variable_0.Code(),
	diagnostics.A_module_cannot_have_multiple_default_exports.Code(),
	diagnostics.Another_export_default_is_here.Code(),
	diagnostics.The_first_export_default_is_here.Code(),
	diagnostics.Identifier_expected_0_is_a_reserved_word_at_the_top_level_of_a_module.Code(),
	diagnostics.Identifier_expected_0_is_a_reserved_word_in_strict_mode_Modules_are_automatically_in_strict_mode.Code(),
	diagnostics.Identifier_expected_0_is_a_reserved_word_that_cannot_be_used_here.Code(),
	diagnostics.X_constructor_is_a_reserved_word.Code(),
	diagnostics.X_delete_cannot_be_called_on_an_identifier_in_strict_mode.Code(),
	diagnostics.Code_contained_in_a_class_is_evaluated_in_JavaScript_s_strict_mode_which_does_not_allow_this_use_of_0_For_more_information_see_https_Colon_Slash_Slashdeveloper_mozilla_org_Slashen_US_Slashdocs_SlashWeb_SlashJavaScript_SlashReference_SlashStrict_mode.Code(),
	diagnostics.Invalid_use_of_0_Modules_are_automatically_in_strict_mode.Code(),
	diagnostics.Invalid_use_of_0_in_strict_mode.Code(),
	diagnostics.A_label_is_not_allowed_here.Code(),
	diagnostics.X_with_statements_are_not_allowed_in_strict_mode.Code(),
	// grammar errors
	diagnostics.A_break_statement_can_only_be_used_within_an_enclosing_iteration_or_switch_statement.Code(),
	diagnostics.A_break_statement_can_only_jump_to_a_label_of_an_enclosing_statement.Code(),
	diagnostics.A_class_declaration_without_the_default_modifier_must_have_a_name.Code(),
	diagnostics.A_class_member_cannot_have_the_0_keyword.Code(),
	diagnostics.A_comma_expression_is_not_allowed_in_a_computed_property_name.Code(),
	diagnostics.A_continue_statement_can_only_be_used_within_an_enclosing_iteration_statement.Code(),
	diagnostics.A_continue_statement_can_only_jump_to_a_label_of_an_enclosing_iteration_statement.Code(),
	diagnostics.A_default_clause_cannot_appear_more_than_once_in_a_switch_statement.Code(),
	diagnostics.A_default_export_must_be_at_the_top_level_of_a_file_or_module_declaration.Code(),
	diagnostics.A_definite_assignment_assertion_is_not_permitted_in_this_context.Code(),
	diagnostics.A_destructuring_declaration_must_have_an_initializer.Code(),
	diagnostics.A_get_accessor_cannot_have_parameters.Code(),
	diagnostics.A_rest_element_cannot_contain_a_binding_pattern.Code(),
	diagnostics.A_rest_element_cannot_have_a_property_name.Code(),
	diagnostics.A_rest_element_cannot_have_an_initializer.Code(),
	diagnostics.A_rest_element_must_be_last_in_a_destructuring_pattern.Code(),
	diagnostics.A_rest_parameter_cannot_have_an_initializer.Code(),
	diagnostics.A_rest_parameter_must_be_last_in_a_parameter_list.Code(),
	diagnostics.A_rest_parameter_or_binding_pattern_may_not_have_a_trailing_comma.Code(),
	diagnostics.A_return_statement_cannot_be_used_inside_a_class_static_block.Code(),
	diagnostics.A_set_accessor_cannot_have_rest_parameter.Code(),
	diagnostics.A_set_accessor_must_have_exactly_one_parameter.Code(),
	diagnostics.An_export_declaration_can_only_be_used_at_the_top_level_of_a_module.Code(),
	diagnostics.An_export_declaration_cannot_have_modifiers.Code(),
	diagnostics.An_import_declaration_can_only_be_used_at_the_top_level_of_a_module.Code(),
	diagnostics.An_import_declaration_cannot_have_modifiers.Code(),
	diagnostics.An_object_member_cannot_be_declared_optional.Code(),
	diagnostics.Argument_of_dynamic_import_cannot_be_spread_element.Code(),
	diagnostics.Cannot_assign_to_private_method_0_Private_methods_are_not_writable.Code(),
	diagnostics.Cannot_redeclare_identifier_0_in_catch_clause.Code(),
	diagnostics.Catch_clause_variable_cannot_have_an_initializer.Code(),
	diagnostics.Class_decorators_can_t_be_used_with_static_private_identifier_Consider_removing_the_experimental_decorator.Code(),
	diagnostics.Classes_can_only_extend_a_single_class.Code(),
	diagnostics.Classes_may_not_have_a_field_named_constructor.Code(),
	diagnostics.Did_you_mean_to_use_a_Colon_An_can_only_follow_a_property_name_when_the_containing_object_literal_is_part_of_a_destructuring_pattern.Code(),
	diagnostics.Duplicate_label_0.Code(),
	diagnostics.Dynamic_imports_can_only_accept_a_module_specifier_and_an_optional_set_of_attributes_as_arguments.Code(),
	diagnostics.X_for_await_loops_cannot_be_used_inside_a_class_static_block.Code(),
	diagnostics.JSX_attributes_must_only_be_assigned_a_non_empty_expression.Code(),
	diagnostics.JSX_elements_cannot_have_multiple_attributes_with_the_same_name.Code(),
	diagnostics.JSX_expressions_may_not_use_the_comma_operator_Did_you_mean_to_write_an_array.Code(),
	diagnostics.JSX_property_access_expressions_cannot_include_JSX_namespace_names.Code(),
	diagnostics.Jump_target_cannot_cross_function_boundary.Code(),
	diagnostics.Line_terminator_not_permitted_before_arrow.Code(),
	diagnostics.Modifiers_cannot_appear_here.Code(),
	diagnostics.Only_a_single_variable_declaration_is_allowed_in_a_for_in_statement.Code(),
	diagnostics.Only_a_single_variable_declaration_is_allowed_in_a_for_of_statement.Code(),
	diagnostics.Private_identifiers_are_not_allowed_outside_class_bodies.Code(),
	diagnostics.Private_identifiers_are_only_allowed_in_class_bodies_and_may_only_be_used_as_part_of_a_class_member_declaration_property_access_or_on_the_left_hand_side_of_an_in_expression.Code(),
	diagnostics.Property_0_is_not_accessible_outside_class_1_because_it_has_a_private_identifier.Code(),
	diagnostics.Tagged_template_expressions_are_not_permitted_in_an_optional_chain.Code(),
	diagnostics.The_left_hand_side_of_a_for_of_statement_may_not_be_async.Code(),
	diagnostics.The_variable_declaration_of_a_for_in_statement_cannot_have_an_initializer.Code(),
	diagnostics.The_variable_declaration_of_a_for_of_statement_cannot_have_an_initializer.Code(),
	diagnostics.Trailing_comma_not_allowed.Code(),
	diagnostics.Variable_declaration_list_cannot_be_empty.Code(),
	diagnostics.X_0_and_1_operations_cannot_be_mixed_without_parentheses.Code(),
	diagnostics.X_0_expected.Code(),
	diagnostics.X_0_is_not_a_valid_meta_property_for_keyword_1_Did_you_mean_2.Code(),
	diagnostics.X_0_list_cannot_be_empty.Code(),
	diagnostics.X_0_modifier_already_seen.Code(),
	diagnostics.X_0_modifier_cannot_appear_on_a_constructor_declaration.Code(),
	diagnostics.X_0_modifier_cannot_appear_on_a_module_or_namespace_element.Code(),
	diagnostics.X_0_modifier_cannot_appear_on_a_parameter.Code(),
	diagnostics.X_0_modifier_cannot_appear_on_class_elements_of_this_kind.Code(),
	diagnostics.X_0_modifier_cannot_be_used_here.Code(),
	diagnostics.X_0_modifier_must_precede_1_modifier.Code(),
	diagnostics.X_0_declarations_can_only_be_declared_inside_a_block.Code(),
	diagnostics.X_0_declarations_must_be_initialized.Code(),
	diagnostics.X_extends_clause_already_seen.Code(),
	diagnostics.X_let_is_not_allowed_to_be_used_as_a_name_in_let_or_const_declarations.Code(),
	diagnostics.Class_constructor_may_not_be_a_generator.Code(),
	diagnostics.Class_constructor_may_not_be_an_accessor.Code(),
	diagnostics.X_await_expressions_are_only_allowed_within_async_functions_and_at_the_top_levels_of_modules.Code(),
	diagnostics.X_await_using_statements_are_only_allowed_within_async_functions_and_at_the_top_levels_of_modules.Code(),
	diagnostics.Private_field_0_must_be_declared_in_an_enclosing_class.Code(),
	// Type errors
	diagnostics.This_condition_will_always_return_0_since_JavaScript_compares_objects_by_reference_not_value.Code(),
)
