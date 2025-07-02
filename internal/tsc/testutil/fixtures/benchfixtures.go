package fixtures

import (
	"path/filepath"

	"github.com/yasufadhili/jawt/internal/tsc/repo"
	"github.com/yasufadhili/jawt/internal/tsc/testutil/filefixture"
)

var BenchFixtures = []filefixture.Fixture{
	filefixture.FromString("empty.ts", "empty.ts", ""),
	filefixture.FromFile("checker.ts", filepath.Join(repo.TypeScriptSubmodulePath, "src/compiler/checker.ts")),
	filefixture.FromFile("dom.generated.d.ts", filepath.Join(repo.TypeScriptSubmodulePath, "src/lib/dom.generated.d.ts")),
	filefixture.FromFile("Herebyfile.mjs", filepath.Join(repo.TypeScriptSubmodulePath, "Herebyfile.mjs")),
	filefixture.FromFile("jsxComplexSignatureHasApplicabilityError.tsx", filepath.Join(repo.TypeScriptSubmodulePath, "tests/cases/compiler/jsxComplexSignatureHasApplicabilityError.tsx")),
}
