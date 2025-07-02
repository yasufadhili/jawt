package fourslash_test

import (
	"testing"

	"github.com/yasufadhili/jawt/internal/tsc/fourslash"
	"github.com/yasufadhili/jawt/internal/tsc/ls"
	"github.com/yasufadhili/jawt/internal/tsc/lsp/lsproto"
	"github.com/yasufadhili/jawt/internal/tsc/testutil"
)

func TestBasicInterfaceMembers(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `export {};
interface Point {
	x: number;
	y: number;
}
declare const p: Point;
p./*a*/`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "a", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:      "x",
					Kind:       ptrTo(lsproto.CompletionItemKindField),
					SortText:   ptrTo(string(ls.SortTextLocationPriority)),
					InsertText: ptrTo(".x"),
					FilterText: ptrTo(".x"),
					TextEdit: &lsproto.TextEditOrInsertReplaceEdit{
						TextEdit: &lsproto.TextEdit{
							NewText: ".x",
							Range: lsproto.Range{
								Start: lsproto.Position{Line: 6, Character: 1},
								End:   lsproto.Position{Line: 6, Character: 2},
							},
						},
					},
				},
				"y",
			},
		},
	})
}
