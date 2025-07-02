package tstransforms_test

import (
	"testing"

	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/binder"
	"github.com/yasufadhili/jawt/internal/tsc/checker"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/module"
	"github.com/yasufadhili/jawt/internal/tsc/modulespecifiers"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/testutil/emittestutil"
	"github.com/yasufadhili/jawt/internal/tsc/testutil/parsetestutil"
	"github.com/yasufadhili/jawt/internal/tsc/transformers/tstransforms"
	"github.com/yasufadhili/jawt/internal/tsc/tsoptions"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
)

type fakeProgram struct {
	singleThreaded                 bool
	compilerOptions                *core.CompilerOptions
	files                          []*ast.SourceFile
	getEmitModuleFormatOfFile      func(sourceFile ast.HasFileName) core.ModuleKind
	getImpliedNodeFormatForEmit    func(sourceFile ast.HasFileName) core.ModuleKind
	getResolvedModule              func(currentSourceFile ast.HasFileName, moduleReference string, mode core.ResolutionMode) *module.ResolvedModule
	getSourceFile                  func(FileName string) *ast.SourceFile
	getSourceFileForResolvedModule func(FileName string) *ast.SourceFile
}

// GetRedirectForResolution implements checker.Program.
func (p *fakeProgram) GetRedirectForResolution(file ast.HasFileName) *tsoptions.ParsedCommandLine {
	panic("unimplemented")
}

// SourceFileMayBeEmitted implements checker.Program.
func (p *fakeProgram) SourceFileMayBeEmitted(sourceFile *ast.SourceFile, forceDtsEmit bool) bool {
	panic("unimplemented")
}

// GetEmitSyntaxForUsageLocation implements checker.Program.
func (p *fakeProgram) GetEmitSyntaxForUsageLocation(sourceFile ast.HasFileName, usageLocation *ast.StringLiteralLike) core.ResolutionMode {
	panic("unimplemented")
}

// CommonSourceDirectory implements checker.Program.
func (p *fakeProgram) CommonSourceDirectory() string {
	panic("unimplemented")
}

func (p *fakeProgram) GetResolvedModuleFromModuleSpecifier(file ast.HasFileName, moduleSpecifier *ast.StringLiteralLike) *module.ResolvedModule {
	panic("unimplemented")
}

func (p *fakeProgram) FileExists(path string) bool {
	return false
}

func (p *fakeProgram) GetCurrentDirectory() string {
	return ""
}

func (p *fakeProgram) GetGlobalTypingsCacheLocation() string {
	return ""
}

func (p *fakeProgram) GetNearestAncestorDirectoryWithPackageJson(dirname string) string {
	return ""
}

func (p *fakeProgram) GetPackageJsonInfo(pkgJsonPath string) modulespecifiers.PackageJsonInfo {
	return nil
}

func (p *fakeProgram) GetRedirectTargets(path tspath.Path) []string {
	return nil
}

func (p *fakeProgram) GetOutputAndProjectReference(path tspath.Path) *tsoptions.OutputDtsAndProjectReference {
	return nil
}

func (p *fakeProgram) IsSourceFromProjectReference(path tspath.Path) bool {
	return false
}

func (p *fakeProgram) GetSourceAndProjectReference(path tspath.Path) *tsoptions.SourceAndProjectReference {
	return nil
}

func (p *fakeProgram) UseCaseSensitiveFileNames() bool {
	return true
}

func (p *fakeProgram) Options() *core.CompilerOptions {
	return p.compilerOptions
}

func (p *fakeProgram) SourceFiles() []*ast.SourceFile {
	return p.files
}

func (p *fakeProgram) BindSourceFiles() {
	wg := core.NewWorkGroup(p.singleThreaded)
	for _, file := range p.files {
		if !file.IsBound() {
			wg.Queue(func() {
				binder.BindSourceFile(file)
			})
		}
	}
	wg.RunAndWait()
}

func (p *fakeProgram) GetEmitModuleFormatOfFile(sourceFile ast.HasFileName) core.ModuleKind {
	return p.getEmitModuleFormatOfFile(sourceFile)
}

func (p *fakeProgram) GetImpliedNodeFormatForEmit(sourceFile ast.HasFileName) core.ModuleKind {
	return p.getImpliedNodeFormatForEmit(sourceFile)
}

func (p *fakeProgram) GetDefaultResolutionModeForFile(sourceFile ast.HasFileName) core.ResolutionMode {
	return p.getEmitModuleFormatOfFile(sourceFile)
}

func (p *fakeProgram) GetModeForUsageLocation(sourceFile ast.HasFileName, location *ast.Node) core.ResolutionMode {
	return p.getEmitModuleFormatOfFile(sourceFile)
}

func (p *fakeProgram) GetResolvedModule(currentSourceFile ast.HasFileName, moduleReference string, mode core.ResolutionMode) *module.ResolvedModule {
	return p.getResolvedModule(currentSourceFile, moduleReference, mode)
}

func (p *fakeProgram) GetSourceFile(FileName string) *ast.SourceFile {
	return p.getSourceFile(FileName)
}

func (p *fakeProgram) GetSourceFileForResolvedModule(FileName string) *ast.SourceFile {
	return p.getSourceFileForResolvedModule(FileName)
}

func (p *fakeProgram) GetSourceFileMetaData(path tspath.Path) ast.SourceFileMetaData {
	return ast.SourceFileMetaData{}
}

func (p *fakeProgram) GetImportHelpersImportSpecifier(path tspath.Path) *ast.Node {
	return nil
}

func (p *fakeProgram) GetJSXRuntimeImportSpecifier(path tspath.Path) (moduleReference string, specifier *ast.Node) {
	return "", nil
}

func (p *fakeProgram) GetResolvedModules() map[tspath.Path]module.ModeAwareCache[*module.ResolvedModule] {
	panic("unimplemented")
}

func (p *fakeProgram) IsSourceFileDefaultLibrary(path tspath.Path) bool {
	return false
}

func TestImportElision(t *testing.T) {
	t.Parallel()
	data := []struct {
		title  string
		input  string
		output string
		other  string
		jsx    bool
	}{
		{title: "ImportEquals#1", input: "import x = require(\"other\"); x;", output: "import x = require(\"other\");\nx;"},
		{title: "ImportEquals#2", input: "import x = require(\"other\");", output: ""},
		{title: "ImportDeclaration#1", input: `import "m";`, output: `import "m";`},
		{title: "ImportDeclaration#2", input: "import * as x from \"other\"; x;", output: "import * as x from \"other\";\nx;"},
		{title: "ImportDeclaration#3", input: "import x from \"other\"; x;", output: "import x from \"other\";\nx;"},
		{title: "ImportDeclaration#4", input: "import { x } from \"other\"; x;", output: "import { x } from \"other\";\nx;"},
		{title: "ImportDeclaration#5", input: "import * as x from \"other\";", output: ""},
		{title: "ImportDeclaration#6", input: "import x from \"other\";", output: ""},
		{title: "ImportDeclaration#7", input: "import { x } from \"other\";", output: ""},
		{title: "ExportDeclaration#1", input: "export * from \"other\";", other: "export let x;", output: "export * from \"other\";"},
		{title: "ExportDeclaration#2", input: "export * as x from \"other\";", other: "export let x;", output: "export * as x from \"other\";"},
		{title: "ExportDeclaration#3", input: "export * from \"other\";", other: "export let x;", output: "export * from \"other\";"},
		{title: "ExportDeclaration#4", input: "export * as x from \"other\";", other: "export let x;", output: "export * as x from \"other\";"},
		{title: "ExportDeclaration#5", input: "export { x } from \"other\";", other: "export let x;", output: "export { x } from \"other\";"},
		{title: "ExportDeclaration#6", input: "export { x } from \"other\";", other: "export type x = any;", output: ""},
		{title: "ExportDeclaration#7", input: "export { x }; let x;", output: "export { x };\nlet x;"},
		{title: "ExportDeclaration#8", input: "export { x }; type x = any;", output: ""},
		{title: "ExportDeclaration#9", input: "import { x } from \"other\"; export { x };", other: "export type x = any;", output: ""},
		{title: "ExportAssignment#1", input: "let x; export default x;", output: "let x;\nexport default x;"},
		{title: "ExportAssignment#2", input: "type x = any; export default x;", output: ""},
	}

	for _, rec := range data {
		t.Run(rec.title, func(t *testing.T) {
			t.Parallel()

			file := parsetestutil.ParseTypeScript(rec.input, rec.jsx)
			parsetestutil.CheckDiagnostics(t, file)
			files := []*ast.SourceFile{file}

			var other *ast.SourceFile
			if len(rec.other) > 0 {
				other = parsetestutil.ParseTypeScript(rec.other, rec.jsx)
				parsetestutil.CheckDiagnostics(t, other)
				files = append(files, other)
			}

			compilerOptions := &core.CompilerOptions{}

			c := checker.NewChecker(&fakeProgram{
				singleThreaded:  true,
				compilerOptions: compilerOptions,
				files:           files,
				getEmitModuleFormatOfFile: func(sourceFile ast.HasFileName) core.ModuleKind {
					return core.ModuleKindESNext
				},
				getImpliedNodeFormatForEmit: func(sourceFile ast.HasFileName) core.ModuleKind {
					return core.ModuleKindESNext
				},
				getSourceFile: func(fileName string) *ast.SourceFile {
					if fileName == "other.ts" {
						return other
					}
					return nil
				},
				getSourceFileForResolvedModule: func(fileName string) *ast.SourceFile {
					if fileName == "other.ts" {
						return other
					}
					return nil
				},
				getResolvedModule: func(currentSourceFile ast.HasFileName, moduleReference string, mode core.ResolutionMode) *module.ResolvedModule {
					if currentSourceFile == file && moduleReference == "other" {
						return &module.ResolvedModule{
							ResolvedFileName: "other.ts",
							Extension:        tspath.ExtensionTs,
						}
					}
					return nil
				},
			})

			emitResolver := c.GetEmitResolver()
			emitResolver.MarkLinkedReferencesRecursively(file)

			emitContext := printer.NewEmitContext()
			file = tstransforms.NewTypeEraserTransformer(emitContext, compilerOptions).TransformSourceFile(file)
			file = tstransforms.NewImportElisionTransformer(emitContext, compilerOptions, emitResolver).TransformSourceFile(file)
			emittestutil.CheckEmit(t, nil, file, rec.output)
		})
	}
}
