package emitter

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/yasufadhili/jawt/internal/ast"
	"github.com/yasufadhili/jawt/internal/core"
)

// PageResult holds the emitted HTML content and its intended output path.
type PageResult struct {
	Content  string
	FilePath string
}

// EmitPage generates the HTML content for a JML page and determines its output path.
func EmitPage(ctx *core.JawtContext, doc *ast.Document) (*PageResult, error) {
	if doc.DocType != ast.DocTypePage {
		return nil, fmt.Errorf("document is not a page: %s", doc.Name.Name)
	}

	var htmlBuilder strings.Builder
	htmlBuilder.WriteString("<!DOCTYPE html>\n")
	htmlBuilder.WriteString("<html>\n")
	htmlBuilder.WriteString("<head>\n")

	// Extract page metadata from AST (e.g., title, description, favicon)
	pageTitle := doc.Name.Name // Default to document name
	// TODO: Extract other page metadata (description, keywords, author, viewport, favicon) from Page AST if available
	// This would require a specific Page node in the AST with these fields, or a way to extract them from the body.

	htmlBuilder.WriteString(fmt.Sprintf("\t<title>%s</title>\n", pageTitle))

	// Collect imported components and generate script tags
	var importedComponents []string
	for _, stmt := range doc.Body {
		if importDecl, ok := stmt.(*ast.ImportDeclaration); ok {
			// All imports in pages are components for now
			// "Pages can only import components"
			for _, spec := range importDecl.Specifiers {
				if importSpec, ok := spec.(*ast.ImportSpecifier); ok {
					importedComponents = append(importedComponents, importSpec.Local.Name)
				} else if importDefaultSpec, ok := spec.(*ast.ImportDefaultSpecifier); ok {
					importedComponents = append(importedComponents, importDefaultSpec.Local.Name)
				} else if importNamespaceSpec, ok := spec.(*ast.ImportNamespaceSpecifier); ok {
					importedComponents = append(importedComponents, importNamespaceSpec.Local.Name)
				}
			}
		}
	}

	// Generate script tags for imported components
	for _, compName := range importedComponents {
		// Components are emitted to ComponentsOutputDir (e.g., .jawt/build/components/)
		// The HTML src attribute needs to be relative to the web root (BuildDir).
		// Example: BuildDir = /project/.jawt/build
		// ComponentsOutputDir = /project/.jawt/build/components
		// Component JS file = /project/.jawt/build/components/MyComponent.js
		// HTML src should be /components/my-component.js

		// Construct the full absolute path to the component's JS file in the build directory
		absoluteCompJSPath := filepath.Join(ctx.Paths.ComponentsOutputDir, sanitiseComponentName(compName)+".js")

		// Get the path relative to the BuildDir (which serves as the web root)
		relCompPath, err := filepath.Rel(ctx.Paths.BuildDir, absoluteCompJSPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get relative path for component %s: %w", compName, err)
		}
		// Prepend / to make it absolute from the web root
		htmlBuilder.WriteString(fmt.Sprintf("\t<script type=\"module\" src=\"/ %s\"></script>\n", relCompPath))
	}

	htmlBuilder.WriteString("</head>\n")
	htmlBuilder.WriteString("<body>\n")

	// Emit the single root component
	// For now, assume the first ComponentElement in the body is the root.
	// This needs to be more robust and handle the actual rendering of the AST.
	var rootComponent *ast.ComponentElement
	for _, stmt := range doc.Body {
		if compElem, ok := stmt.(*ast.ComponentElement); ok {
			rootComponent = compElem
			break
		}
	}

	if rootComponent != nil {
		// This is a simplification. The proper AST visitor will render the component
		// and its properties. For now, we just emit a custom element tag.
		htmlBuilder.WriteString(fmt.Sprintf("\t<%s></%s>\n", strings.ToLower(sanitiseComponentName(rootComponent.Tag.Name)), strings.ToLower(sanitiseComponentName(rootComponent.Tag.Name))))
		// TODO: Pass properties to the root component (requires more complex AST traversal and property mapping)
	} else {
		htmlBuilder.WriteString("\t<!-- No root component found for this page -->\n")
	}

	htmlBuilder.WriteString("</body>\n")
	htmlBuilder.WriteString("</html>\n")

	// Determine output file path based on page routing convention
	// Example: app/about/index.jml -> build/about/index.html

	// Get the relative path of the JML file from the app directory
	relPath, err := filepath.Rel(ctx.Paths.AppDir, doc.SourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path for page %s: %w", doc.SourceFile, err)
	}

	// Change extension to .html
	htmlFileName := strings.TrimSuffix(relPath, filepath.Ext(relPath)) + ".html"

	// Construct the full output path in the build directory
	outputFilePath := filepath.Join(ctx.Paths.BuildDir, htmlFileName)

	return &PageResult{
		Content:  htmlBuilder.String(),
		FilePath: outputFilePath,
	}, nil
}

// sanitiseComponentName converts a PascalCase component name to kebab-case for custom element tags.
func sanitiseComponentName(name string) string {
	var sb strings.Builder
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) {
			sb.WriteRune('-')
		}
		sb.WriteRune(unicode.ToLower(r))
	}
	return sb.String()
}
