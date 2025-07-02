package ls

import (
	"github.com/yasufadhili/jawt/internal/tsc/core"
)

type Location struct {
	FileName string
	Range    core.TextRange
}

type JsxAttributeCompletionStyle string

const (
	JsxAttributeCompletionStyleAuto   JsxAttributeCompletionStyle = "auto"
	JsxAttributeCompletionStyleBraces JsxAttributeCompletionStyle = "braces"
	JsxAttributeCompletionStyleNone   JsxAttributeCompletionStyle = "none"
)

type UserPreferences struct {
	// Enables auto-import-style completions on partially-typed import statements. E.g., allows
	// `import write|` to be completed to `import { writeFile } from "fs"`.
	IncludeCompletionsForImportStatements *bool

	// Unless this option is `false`,  member completion lists triggered with `.` will include entries
	// on potentially-null and potentially-undefined values, with insertion text to replace
	// preceding `.` tokens with `?.`.
	IncludeAutomaticOptionalChainCompletions *bool

	// If enabled, completions for class members (e.g. methods and properties) will include
	// a whole declaration for the member.
	// E.g., `class A { f| }` could be completed to `class A { foo(): number {} }`, instead of
	// `class A { foo }`.
	IncludeCompletionsWithClassMemberSnippets *bool

	// If enabled, object literal methods will have a method declaration completion entry in addition
	// to the regular completion entry containing just the method name.
	// E.g., `const objectLiteral: T = { f| }` could be completed to `const objectLiteral: T = { foo(): void {} }`,
	// in addition to `const objectLiteral: T = { foo }`.
	IncludeCompletionsWithObjectLiteralMethodSnippets *bool

	JsxAttributeCompletionStyle *JsxAttributeCompletionStyle
}
