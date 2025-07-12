package ast

import (
	"bytes"
	"strings"
	"testing"
)

// TestPrinter is a comprehensive test that builds a complex AST and
// uses the Printer to verify its structure.
func TestPrinter(t *testing.T) {
	// 1. Construct the AST using factory methods
	pos := Position{Line: 1, Column: 1, File: "test.jml"}

	ast := NewProgram(pos, []*Document{
		NewDocument(pos, DocTypePage, NewIdentifier(pos, "HomePage"), []Statement{
			// import component Card from "./Card.jml"
			NewImportDeclaration(pos, []Node{
				NewImportSpecifier(pos, NewIdentifier(pos, "Card")), // Local: Card, Imported: nil
			}, NewLiteral(pos, "string", "./Card.jml"), false),

			// export default const name = "World"
			NewExportDeclaration(pos, NewVariableDeclaration(pos, "const", []*VariableDeclarator{
				NewVariableDeclarator(pos, NewIdentifier(pos, "name"), NewLiteral(pos, "string", "World")),
			}), true),

			// state count: number = 0
			NewStateDeclaration(pos, NewIdentifier(pos, "count"), NewTypeAnnotation(pos, NewTypeReference(pos, NewIdentifier(pos, "number"), nil)), NewLiteral(pos, "number", "0")),

			// ComponentElement
			NewComponentElement(pos, NewIdentifier(pos, "View"), []ComponentBodyElement{
				// title: `Hello, ${name}`
				NewComponentProperty(pos, NewIdentifier(pos, "title"), NewTemplateLiteral(pos, nil, []Expression{NewIdentifier(pos, "name")})),

				// Card { content: "Welcome!" }
				NewComponentElement(pos, NewIdentifier(pos, "Card"), []ComponentBodyElement{
					NewComponentProperty(pos, NewIdentifier(pos, "content"), NewLiteral(pos, "string", "Welcome!")),
				}),

				// for (item in items) { ... }
				NewForLoop(pos, NewIdentifier(pos, "item"), NewIdentifier(pos, "items"), []ComponentBodyElement{
					NewComponentElement(pos, NewIdentifier(pos, "ListItem"), nil),
				}),

				// if (count > 0) { ... }
				NewIfCondition(pos, NewBinaryExpression(pos, ">", NewIdentifier(pos, "count"), NewLiteral(pos, "number", "0")), []ComponentBodyElement{
					NewComponentElement(pos, NewIdentifier(pos, "Button"), nil),
				}, nil),
			}),
		}, "test.jml"),
	})

	// 2. Print the AST to a buffer
	var buf bytes.Buffer
	printer := NewPrinter(&buf)
	printer.Print(ast)

	// 3. Compare with the expected output
	expected := `
Program
  Document (page HomePage)
    ImportDeclaration from "./Card.jml"
    ExportDeclaration (default )
      VariableDeclaration (const)
        Declarator (name)
          Literal (string: World)
    StateDeclaration
    ComponentElement <View>
      ComponentProperty (title)
        TemplateLiteral
          Identifier (name)
      ComponentElement <Card>
        ComponentProperty (content)
          Literal (string: Welcome!)
      ForLoop (item in ...)
        Source:
          Identifier (items)
        Body:
          ComponentElement <ListItem>
      IfCondition
        Test:
          BinaryExpression (>)
            Identifier (count)
            Literal (number: 0)
        Consequent:
          ComponentElement <Button>
`

	// Normalise whitespace for comparison
	normalize := func(s string) string {
		return strings.Join(strings.Fields(s), " ")
	}

	if normalize(buf.String()) != normalize(expected) {
		t.Errorf("Printer output does not match expected output.\nGot:\n%s\nExpected:\n%s", buf.String(), expected)
	}
}
