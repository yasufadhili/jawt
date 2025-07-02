package fourslash_test

import (
	"testing"

	"github.com/yasufadhili/jawt/internal/tsc/fourslash"
	"github.com/yasufadhili/jawt/internal/tsc/ls"
	"github.com/yasufadhili/jawt/internal/tsc/lsp/lsproto"
	"github.com/yasufadhili/jawt/internal/tsc/testutil"
)

func TestBasicMultifileCompletions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /a.ts
export const foo = { bar: 'baz' };

// @Filename: /b.ts
import { foo } from './a';
const test = foo./*1*/`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "1", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:      "bar",
					Kind:       ptrTo(lsproto.CompletionItemKindField),
					SortText:   ptrTo(string(ls.SortTextLocationPriority)),
					FilterText: ptrTo(".bar"),
					InsertText: ptrTo(".bar"),
					TextEdit: &lsproto.TextEditOrInsertReplaceEdit{
						TextEdit: &lsproto.TextEdit{
							NewText: ".bar",
							Range: lsproto.Range{
								Start: lsproto.Position{Line: 1, Character: 16},
								End:   lsproto.Position{Line: 1, Character: 17},
							},
						},
					},
				},
			},
		},
	})
}
