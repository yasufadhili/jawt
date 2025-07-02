package fourslash_test

import (
	"testing"

	"github.com/yasufadhili/jawt/internal/tsc/fourslash"
	"github.com/yasufadhili/jawt/internal/tsc/testutil"
)

func TestBasicClassElementKeywords(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class C {
	/*a*/
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "a", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: completionClassElementKeywords,
		},
	})
}
