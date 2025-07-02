package ls

import (
	"github.com/yasufadhili/jawt/internal/tsc/astnav"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/lsp/lsproto"
)

func (l *LanguageService) GetExpectedReferenceFromMarker(fileName string, pos int) *lsproto.Location {
	// Temporary testing function--this function only works for markers that are on symbols/names.
	// We won't need this once marker ranges are implemented, or once reference tests are baselined
	_, sourceFile := l.tryGetProgramAndFile(fileName)
	node := astnav.GetTouchingPropertyName(sourceFile, pos)
	return &lsproto.Location{
		Uri:   FileNameToDocumentURI(fileName),
		Range: *l.createLspRangeFromNode(node, sourceFile),
	}
}

func (l *LanguageService) TestProvideReferences(fileName string, pos int) []*lsproto.Location {
	_, sourceFile := l.tryGetProgramAndFile(fileName)
	lsPos := l.converters.PositionToLineAndCharacter(sourceFile, core.TextPos(pos))
	return l.ProvideReferences(&lsproto.ReferenceParams{
		TextDocumentPositionParams: lsproto.TextDocumentPositionParams{
			TextDocument: lsproto.TextDocumentIdentifier{
				Uri: FileNameToDocumentURI(fileName),
			},
			Position: lsPos,
		},
	})
}
