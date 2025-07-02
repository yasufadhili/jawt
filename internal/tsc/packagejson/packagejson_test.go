package packagejson_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	json2 "github.com/go-json-experiment/json"
	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/packagejson"
	"github.com/yasufadhili/jawt/internal/tsc/parser"
	"github.com/yasufadhili/jawt/internal/tsc/repo"
	"github.com/yasufadhili/jawt/internal/tsc/testutil/filefixture"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
)

var packageJsonFixtures = []filefixture.Fixture{
	filefixture.FromFile("package.json", filepath.Join(repo.RootPath, "package.json")),
	filefixture.FromFile("date-fns.json", filepath.Join(repo.TestDataPath, "fixtures", "packagejson", "date-fns.json")),
}

func BenchmarkPackageJSON(b *testing.B) {
	for _, f := range packageJsonFixtures {
		f.SkipIfNotExist(b)
		content := []byte(f.ReadFile(b))
		b.Run("UnmarshalJSON", func(b *testing.B) {
			b.Run(f.Name(), func(b *testing.B) {
				for b.Loop() {
					var p packagejson.Fields
					if err := json.Unmarshal(content, &p); err != nil {
						b.Fatal(err)
					}
				}
			})
		})

		b.Run("UnmarshalJSONV2", func(b *testing.B) {
			b.Run(f.Name(), func(b *testing.B) {
				for b.Loop() {
					var p packagejson.Fields
					if err := json2.Unmarshal(content, &p); err != nil {
						b.Fatal(err)
					}
				}
			})
		})

		b.Run("ParseJSONText", func(b *testing.B) {
			b.Run(f.Name(), func(b *testing.B) {
				fileName := "/" + f.Name()
				for b.Loop() {
					parser.ParseSourceFile(ast.SourceFileParseOptions{
						FileName: fileName,
						Path:     tspath.Path(fileName),
					}, string(content), core.ScriptKindJSON)
				}
			})
		})
	}
}
