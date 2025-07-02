package moduletransforms

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/outputpaths"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
)

func isDeclarationNameOfEnumOrNamespace(emitContext *printer.EmitContext, node *ast.IdentifierNode) bool {
	if original := emitContext.MostOriginal(node); original != nil && original.Parent != nil {
		switch original.Parent.Kind {
		case ast.KindEnumDeclaration, ast.KindModuleDeclaration:
			return original == original.Parent.Name()
		}
	}
	return false
}

func rewriteModuleSpecifier(emitContext *printer.EmitContext, node *ast.Expression, compilerOptions *core.CompilerOptions) *ast.Expression {
	if node == nil || !ast.IsStringLiteral(node) || !core.ShouldRewriteModuleSpecifier(node.Text(), compilerOptions) {
		return node
	}
	updatedText := tspath.ChangeExtension(node.Text(), outputpaths.GetOutputExtension(node.Text(), compilerOptions.Jsx))
	if updatedText != node.Text() {
		updated := emitContext.Factory.NewStringLiteral(updatedText)
		// !!! set quote style
		emitContext.SetOriginal(updated, node)
		emitContext.AssignCommentAndSourceMapRanges(updated, node)
		return updated
	}
	return node
}

func createEmptyImports(factory *printer.NodeFactory) *ast.Statement {
	return factory.NewExportDeclaration(
		nil,   /*modifiers*/
		false, /*isTypeOnly*/
		factory.NewNamedExports(factory.NewNodeList(nil)),
		nil, /*moduleSpecifier*/
		nil, /*attributes*/
	)
}

// Get the name of a target module from an import/export declaration as should be written in the emitted output.
// The emitted output name can be different from the input if:
//  1. The module has a /// <amd-module name="<new name>" />
//  2. --out or --outFile is used, making the name relative to the rootDir
//     3- The containing SourceFile has an entry in renamedDependencies for the import as requested by some module loaders (e.g. System).
//
// Otherwise, a new StringLiteral node representing the module name will be returned.
func getExternalModuleNameLiteral(factory *printer.NodeFactory, importNode *ast.Node /*ImportDeclaration | ExportDeclaration | ImportEqualsDeclaration | ImportCall*/, sourceFile *ast.SourceFile, host any /*EmitHost*/, resolver printer.EmitResolver, compilerOptions *core.CompilerOptions) *ast.StringLiteralNode {
	moduleName := ast.GetExternalModuleName(importNode)
	if moduleName != nil && ast.IsStringLiteral(moduleName) {
		name := tryGetModuleNameFromDeclaration(importNode, host, factory, resolver, compilerOptions)
		if name == nil {
			name = tryRenameExternalModule(factory, moduleName, sourceFile)
		}
		if name == nil {
			name = factory.NewStringLiteral(moduleName.Text())
		}
		return name
	}
	return nil
}

// Get the name of a module as should be written in the emitted output.
// The emitted output name can be different from the input if:
//  1. The module has a /// <amd-module name="<new name>" />
//  2. --out or --outFile is used, making the name relative to the rootDir
//
// Otherwise, a new StringLiteral node representing the module name will be returned.
func tryGetModuleNameFromFile(factory *printer.NodeFactory, file *ast.SourceFile, host any /*EmitHost*/, options *core.CompilerOptions) *ast.StringLiteralNode {
	if file == nil {
		return nil
	}
	// !!!
	// if file.moduleName {
	// 	return factory.createStringLiteral(file.moduleName)
	// }
	if !file.IsDeclarationFile && len(options.OutFile) > 0 {
		return factory.NewStringLiteral(getExternalModuleNameFromPath(host, file.FileName(), "" /*referencePath*/))
	}
	return nil
}

func tryGetModuleNameFromDeclaration(declaration *ast.Node /*ImportEqualsDeclaration | ImportDeclaration | ExportDeclaration | ImportCall*/, host any /*EmitHost*/, factory *printer.NodeFactory, resolver printer.EmitResolver, compilerOptions *core.CompilerOptions) *ast.StringLiteralNode {
	if resolver == nil {
		return nil
	}
	return tryGetModuleNameFromFile(factory, resolver.GetExternalModuleFileFromDeclaration(declaration), host, compilerOptions)
}

// Resolves a local path to a path which is absolute to the base of the emit
func getExternalModuleNameFromPath(host any /*ResolveModuleNameResolutionHost*/, fileName string, referencePath string) string {
	// !!!
	return ""
}

// Some bundlers (SystemJS builder) sometimes want to rename dependencies.
// Here we check if alternative name was provided for a given moduleName and return it if possible.
func tryRenameExternalModule(factory *printer.NodeFactory, moduleName *ast.LiteralExpression, sourceFile *ast.SourceFile) *ast.StringLiteralNode {
	// !!!
	return nil
}

func isFileLevelReservedGeneratedIdentifier(emitContext *printer.EmitContext, name *ast.IdentifierNode) bool {
	info := emitContext.GetAutoGenerateInfo(name)
	return info != nil &&
		info.Flags.IsFileLevel() &&
		info.Flags.IsOptimistic() &&
		info.Flags.IsReservedInNestedScopes()
}

// Used in the module transformer to check if an expression is reasonably without sideeffect,
//
//	and thus better to copy into multiple places rather than to cache in a temporary variable
//	- this is mostly subjective beyond the requirement that the expression not be sideeffecting
func isSimpleCopiableExpression(expression *ast.Expression) bool {
	return ast.IsStringLiteralLike(expression) ||
		ast.IsNumericLiteral(expression) ||
		ast.IsKeywordKind(expression.Kind) ||
		ast.IsIdentifier(expression)
}

// A simple inlinable expression is an expression which can be copied into multiple locations
// without risk of repeating any sideeffects and whose value could not possibly change between
// any such locations
func isSimpleInlineableExpression(expression *ast.Expression) bool {
	return !ast.IsIdentifier(expression) && isSimpleCopiableExpression(expression)
}
