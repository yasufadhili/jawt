package moduletransforms

import (
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/binder"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/printer"
	"github.com/yasufadhili/jawt/internal/tsc/transformers"
)

type ImpliedModuleTransformer struct {
	transformers.Transformer
	compilerOptions           *core.CompilerOptions
	resolver                  binder.ReferenceResolver
	getEmitModuleFormatOfFile func(file ast.HasFileName) core.ModuleKind
	cjsTransformer            *transformers.Transformer
	esmTransformer            *transformers.Transformer
}

func NewImpliedModuleTransformer(emitContext *printer.EmitContext, compilerOptions *core.CompilerOptions, resolver binder.ReferenceResolver, getEmitModuleFormatOfFile func(file ast.HasFileName) core.ModuleKind) *transformers.Transformer {
	if resolver == nil {
		resolver = binder.NewReferenceResolver(compilerOptions, binder.ReferenceResolverHooks{})
	}
	tx := &ImpliedModuleTransformer{compilerOptions: compilerOptions, resolver: resolver, getEmitModuleFormatOfFile: getEmitModuleFormatOfFile}
	return tx.NewTransformer(tx.visit, emitContext)
}

func (tx *ImpliedModuleTransformer) visit(node *ast.Node) *ast.Node {
	switch node.Kind {
	case ast.KindSourceFile:
		node = tx.visitSourceFile(node.AsSourceFile())
	}
	return node
}

func (tx *ImpliedModuleTransformer) visitSourceFile(node *ast.SourceFile) *ast.Node {
	if node.IsDeclarationFile {
		return node.AsNode()
	}

	format := tx.getEmitModuleFormatOfFile(node)

	var transformer *transformers.Transformer
	if format >= core.ModuleKindES2015 {
		if tx.esmTransformer == nil {
			tx.esmTransformer = NewESModuleTransformer(tx.EmitContext(), tx.compilerOptions, tx.resolver, tx.getEmitModuleFormatOfFile)
		}
		transformer = tx.esmTransformer
	} else {
		if tx.cjsTransformer == nil {
			tx.cjsTransformer = NewCommonJSModuleTransformer(tx.EmitContext(), tx.compilerOptions, tx.resolver, tx.getEmitModuleFormatOfFile)
		}
		transformer = tx.cjsTransformer
	}

	return transformer.TransformSourceFile(node).AsNode()
}
