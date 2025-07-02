package fourslash_test

import (
	"testing"

	"github.com/yasufadhili/jawt/internal/tsc/fourslash"
	"github.com/yasufadhili/jawt/internal/tsc/testutil"
)

func TestCompletionListDefaultTypeArgumentPositionTypeOnly(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const foo = "foo";
function test1<T = /*1*/>() {}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "1", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: completionGlobalTypes,
		},
	})
}
