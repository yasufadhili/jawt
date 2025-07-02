package ls_test

import (
	"testing"

	"github.com/yasufadhili/jawt/internal/tsc/bundled"
	"github.com/yasufadhili/jawt/internal/tsc/fourslash"
	"github.com/yasufadhili/jawt/internal/tsc/ls"
	"github.com/yasufadhili/jawt/internal/tsc/lsp/lsproto"
	"github.com/yasufadhili/jawt/internal/tsc/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func TestDefinition(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		// Without embedding, we'd need to read all of the lib files out from disk into the MapFS.
		// Just skip this for now.
		t.Skip("bundled files are not embedded")
	}

	testCases := []struct {
		title    string
		input    string
		expected map[string]lsproto.Definition
	}{
		{
			title: "localFunction",
			input: `
// @filename: index.ts
function localFunction() { }
/*localFunction*/localFunction();`,
			expected: map[string]lsproto.Definition{
				"localFunction": {
					Locations: &[]lsproto.Location{{
						Uri:   ls.FileNameToDocumentURI("/index.ts"),
						Range: lsproto.Range{Start: lsproto.Position{Character: 9}, End: lsproto.Position{Character: 22}},
					}},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			t.Parallel()
			runDefinitionTest(t, testCase.input, testCase.expected)
		})
	}
}

func runDefinitionTest(t *testing.T, input string, expected map[string]lsproto.Definition) {
	testData := fourslash.ParseTestData(t, input, "/mainFile.ts")
	file := testData.Files[0].FileName()
	markerPositions := testData.MarkerPositions
	ctx := projecttestutil.WithRequestID(t.Context())
	languageService, done := createLanguageService(ctx, file, map[string]any{
		file: testData.Files[0].Content,
	})
	defer done()

	for markerName, expectedResult := range expected {
		marker, ok := markerPositions[markerName]
		if !ok {
			t.Fatalf("No marker found for '%s'", markerName)
		}
		locations, err := languageService.ProvideDefinition(
			ctx,
			ls.FileNameToDocumentURI(file),
			marker.LSPosition)
		assert.NilError(t, err)
		assert.DeepEqual(t, *locations, expectedResult)
	}
}
