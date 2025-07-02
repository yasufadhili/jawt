package tsoptions_test

import (
	"slices"
	"testing"

	"github.com/yasufadhili/jawt/internal/tsc/tsoptions"
	"github.com/yasufadhili/jawt/internal/tsc/tsoptions/tsoptionstest"
	"github.com/yasufadhili/jawt/internal/tsc/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestParsedCommandLine(t *testing.T) {
	t.Parallel()
	t.Run("MatchesFileName", func(t *testing.T) {
		t.Parallel()

		noFiles := map[string]string{}
		noFilesFS := vfstest.FromMap(noFiles, true)

		files := map[string]string{
			"/dev/a.ts":         "",
			"/dev/a.d.ts":       "",
			"/dev/a.js":         "",
			"/dev/b.ts":         "",
			"/dev/b.js":         "",
			"/dev/c.d.ts":       "",
			"/dev/z/a.ts":       "",
			"/dev/z/abz.ts":     "",
			"/dev/z/aba.ts":     "",
			"/dev/z/b.ts":       "",
			"/dev/z/bbz.ts":     "",
			"/dev/z/bba.ts":     "",
			"/dev/x/a.ts":       "",
			"/dev/x/aa.ts":      "",
			"/dev/x/b.ts":       "",
			"/dev/x/y/a.ts":     "",
			"/dev/x/y/b.ts":     "",
			"/dev/js/a.js":      "",
			"/dev/js/b.js":      "",
			"/dev/js/d.min.js":  "",
			"/dev/js/ab.min.js": "",
			"/ext/ext.ts":       "",
			"/ext/b/a..b.ts":    "",
		}

		assertMatches := func(t *testing.T, parsedCommandLine *tsoptions.ParsedCommandLine, files map[string]string, matches []string) {
			t.Helper()
			for fileName := range files {
				actual := parsedCommandLine.MatchesFileName(fileName)
				expected := slices.Contains(matches, fileName)
				assert.Equal(t, actual, expected, "fileName: %s", fileName)
			}
			for _, fileName := range matches {
				if _, ok := files[fileName]; !ok {
					actual := parsedCommandLine.MatchesFileName(fileName)
					assert.Equal(t, actual, true, "fileName: %s", fileName)
				}
			}
		}

		t.Run("with literal file list", func(t *testing.T) {
			t.Parallel()
			t.Run("without exclude", func(t *testing.T) {
				t.Parallel()
				parsedCommandLine := tsoptionstest.GetParsedCommandLine(
					t,
					`{
						"files": [
							"a.ts",
							"b.ts"
						]
					}`,
					files,
					"/dev",
					/*useCaseSensitiveFileNames*/ true,
				)

				assertMatches(t, parsedCommandLine, files, []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})
			})

			t.Run("are not removed due to excludes", func(t *testing.T) {
				t.Parallel()
				parsedCommandLine := tsoptionstest.GetParsedCommandLine(
					t,
					`{
						"files": [
							"a.ts",
							"b.ts"
						],
						"exclude": [
							"b.ts"
						]
					}`,
					files,
					"/dev",
					/*useCaseSensitiveFileNames*/ true,
				)

				assertMatches(t, parsedCommandLine, files, []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})

				emptyParsedCommandLine := tsoptions.ReloadFileNamesOfParsedCommandLine(parsedCommandLine, noFilesFS)
				assertMatches(t, emptyParsedCommandLine, noFiles, []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})
			})
		})

		t.Run("with literal include list", func(t *testing.T) {
			t.Parallel()
			t.Run("without exclude", func(t *testing.T) {
				t.Parallel()
				parsedCommandLine := tsoptionstest.GetParsedCommandLine(
					t,
					`{
						"include": [
							"a.ts",
							"b.ts"
						]
					}`,
					files,
					"/dev",
					/*useCaseSensitiveFileNames*/ true,
				)

				assertMatches(t, parsedCommandLine, files, []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})

				emptyParsedCommandLine := tsoptions.ReloadFileNamesOfParsedCommandLine(parsedCommandLine, noFilesFS)
				assertMatches(t, emptyParsedCommandLine, noFiles, []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})
			})

			t.Run("with non .ts file extensions", func(t *testing.T) {
				t.Parallel()
				parsedCommandLine := tsoptionstest.GetParsedCommandLine(
					t,
					`{
						"include": [
							"a.js",
							"b.js"
						]
					}`,
					files,
					"/dev",
					/*useCaseSensitiveFileNames*/ true,
				)

				assertMatches(t, parsedCommandLine, files, []string{})

				emptyParsedCommandLine := tsoptions.ReloadFileNamesOfParsedCommandLine(parsedCommandLine, noFilesFS)
				assertMatches(t, emptyParsedCommandLine, noFiles, []string{})
			})

			t.Run("with literal excludes", func(t *testing.T) {
				t.Parallel()
				parsedCommandLine := tsoptionstest.GetParsedCommandLine(
					t,
					`{
						"include": [
							"a.ts",
							"b.ts"
						],
						"exclude": [
							"b.ts"
						]
					}`,
					files,
					"/dev",
					/*useCaseSensitiveFileNames*/ true,
				)

				assertMatches(t, parsedCommandLine, files, []string{
					"/dev/a.ts",
				})

				emptyParsedCommandLine := tsoptions.ReloadFileNamesOfParsedCommandLine(parsedCommandLine, noFilesFS)
				assertMatches(t, emptyParsedCommandLine, noFiles, []string{
					"/dev/a.ts",
				})
			})

			t.Run("with wildcard excludes", func(t *testing.T) {
				t.Parallel()
				parsedCommandLine := tsoptionstest.GetParsedCommandLine(
					t,
					`{
						"include": [
							"a.ts",
							"b.ts",
							"z/a.ts",
							"z/abz.ts",
							"z/aba.ts",
							"x/b.ts"
						],
						"exclude": [
							"*.ts",
							"z/??z.ts",
							"*/b.ts"
						]
					}`,
					files,
					"/dev",
					/*useCaseSensitiveFileNames*/ true,
				)

				assertMatches(t, parsedCommandLine, files, []string{
					"/dev/z/a.ts",
					"/dev/z/aba.ts",
				})

				emptyParsedCommandLine := tsoptions.ReloadFileNamesOfParsedCommandLine(parsedCommandLine, noFilesFS)
				assertMatches(t, emptyParsedCommandLine, noFiles, []string{
					"/dev/z/a.ts",
					"/dev/z/aba.ts",
				})
			})

			t.Run("with wildcard include list", func(t *testing.T) {
				t.Parallel()

				t.Run("star matches only ts files", func(t *testing.T) {
					t.Parallel()
					parsedCommandLine := tsoptionstest.GetParsedCommandLine(
						t,
						`{
							"include": [
								"*"
							]
						}`,
						files,
						"/dev",
						/*useCaseSensitiveFileNames*/ true,
					)

					assertMatches(t, parsedCommandLine, files, []string{
						"/dev/a.ts",
						"/dev/b.ts",
						"/dev/c.d.ts",
					})

					// a.d.ts matches if a.ts is not already included
					emptyParsedCommandLine := tsoptions.ReloadFileNamesOfParsedCommandLine(parsedCommandLine, noFilesFS)
					assertMatches(t, emptyParsedCommandLine, noFiles, []string{
						"/dev/a.ts",
						"/dev/a.d.ts",
						"/dev/b.ts",
						"/dev/c.d.ts",
					})
				})

				t.Run("question matches only a single character", func(t *testing.T) {
					t.Parallel()
					parsedCommandLine := tsoptionstest.GetParsedCommandLine(
						t,
						`{
							"include": [
								"x/?.ts"
							]
						}`,
						files,
						"/dev",
						/*useCaseSensitiveFileNames*/ true,
					)

					assertMatches(t, parsedCommandLine, files, []string{
						"/dev/x/a.ts",
						"/dev/x/b.ts",
					})

					emptyParsedCommandLine := tsoptions.ReloadFileNamesOfParsedCommandLine(parsedCommandLine, noFilesFS)
					assertMatches(t, emptyParsedCommandLine, noFiles, []string{
						"/dev/x/a.ts",
						"/dev/x/b.ts",
					})
				})

				t.Run("exclude .js files when allowJs=false", func(t *testing.T) {
					t.Parallel()
					parsedCommandLine := tsoptionstest.GetParsedCommandLine(
						t,
						`{
							"include": [
								"js/*"
							]
						}`,
						files,
						"/dev",
						/*useCaseSensitiveFileNames*/ true,
					)

					assertMatches(t, parsedCommandLine, files, []string{})

					emptyParsedCommandLine := tsoptions.ReloadFileNamesOfParsedCommandLine(parsedCommandLine, noFilesFS)
					assertMatches(t, emptyParsedCommandLine, noFiles, []string{})
				})

				t.Run("include .js files when allowJs=true", func(t *testing.T) {
					t.Parallel()
					parsedCommandLine := tsoptionstest.GetParsedCommandLine(
						t,
						`{
							"compilerOptions": {
								"allowJs": true
							},
							"include": [
								"js/*"
							]
						}`,
						files,
						"/dev",
						/*useCaseSensitiveFileNames*/ true,
					)

					assertMatches(t, parsedCommandLine, files, []string{
						"/dev/js/a.js",
						"/dev/js/b.js",
					})

					emptyParsedCommandLine := tsoptions.ReloadFileNamesOfParsedCommandLine(parsedCommandLine, noFilesFS)
					assertMatches(t, emptyParsedCommandLine, noFiles, []string{
						"/dev/js/a.js",
						"/dev/js/b.js",
					})
				})

				t.Run("include explicitly listed .min.js files when allowJs=true", func(t *testing.T) {
					t.Parallel()
					parsedCommandLine := tsoptionstest.GetParsedCommandLine(
						t,
						`{
							"compilerOptions": {
								"allowJs": true
							},
							"include": [
								"js/*.min.js"
							]
						}`,
						files,
						"/dev",
						/*useCaseSensitiveFileNames*/ true,
					)

					assertMatches(t, parsedCommandLine, files, []string{
						"/dev/js/d.min.js",
						"/dev/js/ab.min.js",
					})

					emptyParsedCommandLine := tsoptions.ReloadFileNamesOfParsedCommandLine(parsedCommandLine, noFilesFS)
					assertMatches(t, emptyParsedCommandLine, noFiles, []string{
						"/dev/js/d.min.js",
						"/dev/js/ab.min.js",
					})
				})
			})
		})
	})
}
