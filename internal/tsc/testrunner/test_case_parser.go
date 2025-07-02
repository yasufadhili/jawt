package testrunner

import (
	"regexp"
	"slices"
	"strings"

	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/core"
	"github.com/yasufadhili/jawt/internal/tsc/parser"
	"github.com/yasufadhili/jawt/internal/tsc/scanner"
	"github.com/yasufadhili/jawt/internal/tsc/testutil/harnessutil"
	"github.com/yasufadhili/jawt/internal/tsc/tsoptions"
	"github.com/yasufadhili/jawt/internal/tsc/tsoptions/tsoptionstest"
	"github.com/yasufadhili/jawt/internal/tsc/tspath"
)

var lineDelimiter = regexp.MustCompile("\r?\n")

// This maps a compiler setting to its value as written in the test file. For example, if a test file contains:
//
//	// @target: esnext, es2015
//
// Then the map will map "target" to "esnext, es2015"
type rawCompilerSettings map[string]string

// All the necessary information to turn a multi file test into useful units for later compilation
type testUnit struct {
	content string
	name    string
}

type testCaseContent struct {
	testUnitData         []*testUnit
	tsConfig             *tsoptions.ParsedCommandLine
	tsConfigFileUnitData *testUnit
	symlinks             map[string]string
}

// Regex for parsing options in the format "@Alpha: Value of any sort"
var optionRegex = regexp.MustCompile(`(?m)^\/{2}\s*@(\w+)\s*:\s*([^\r\n]*)`)

// Regex for parsing @link option
var linkRegex = regexp.MustCompile(`(?m)^\/{2}\s*@link\s*:\s*([^\r\n]*)\s*->\s*([^\r\n]*)`)

// File-specific directives used by fourslash tests
var fourslashDirectives = []string{"emitthisfile"}

// Given a test file containing // @FileName directives,
// return an array of named units of code to be added to an existing compiler instance.
func makeUnitsFromTest(code string, fileName string) testCaseContent {
	testUnits, symlinks, currentDirectory, _, _ := ParseTestFilesAndSymlinks(
		code,
		fileName,
		func(filename string, content string, fileOptions map[string]string) (*testUnit, error) {
			return &testUnit{content: content, name: filename}, nil
		},
	)
	if currentDirectory == "" {
		currentDirectory = srcFolder
	}

	// unit tests always list files explicitly
	allFiles := make(map[string]string)
	for _, data := range testUnits {
		allFiles[tspath.GetNormalizedAbsolutePath(data.name, currentDirectory)] = data.content
	}
	parseConfigHost := tsoptionstest.NewVFSParseConfigHost(allFiles, currentDirectory, true /*useCaseSensitiveFileNames*/)

	// check if project has tsconfig.json in the list of files
	var tsConfig *tsoptions.ParsedCommandLine
	var tsConfigFileUnitData *testUnit
	for i, data := range testUnits {
		if harnessutil.GetConfigNameFromFileName(data.name) != "" {
			configFileName := tspath.GetNormalizedAbsolutePath(data.name, currentDirectory)
			path := tspath.ToPath(data.name, parseConfigHost.GetCurrentDirectory(), parseConfigHost.Vfs.UseCaseSensitiveFileNames())
			configJson := parser.ParseSourceFile(ast.SourceFileParseOptions{
				FileName: configFileName,
				Path:     path,
			}, data.content, core.ScriptKindJSON)
			tsConfigSourceFile := &tsoptions.TsConfigSourceFile{
				SourceFile: configJson,
			}
			configDir := tspath.GetDirectoryPath(configFileName)
			tsConfig = tsoptions.ParseJsonSourceFileConfigFileContent(
				tsConfigSourceFile,
				parseConfigHost,
				configDir,
				nil, /*existingOptions*/
				configFileName,
				nil, /*resolutionStack*/
				nil, /*extraFileExtensions*/
				nil /*extendedConfigCache*/)
			tsConfigFileUnitData = data

			// delete tsconfig file entry from the list
			testUnits = slices.Delete(testUnits, i, i+1)
			break
		}
	}

	return testCaseContent{
		testUnitData:         testUnits,
		tsConfig:             tsConfig,
		tsConfigFileUnitData: tsConfigFileUnitData,
		symlinks:             symlinks,
	}
}

// Given a test file containing // @FileName and // @symlink directives,
// return an array of named units of code to be added to an existing compiler instance,
// along with a map of symlinks and the current directory.
func ParseTestFilesAndSymlinks[T any](
	code string,
	fileName string,
	parseFile func(filename string, content string, fileOptions map[string]string) (T, error),
) (units []T, symlinks map[string]string, currentDir string, globalOptions map[string]string, e error) {
	// List of all the subfiles we've parsed out
	var testUnits []T

	lines := lineDelimiter.Split(code, -1)

	// Stuff related to the subfile we're parsing
	var currentFileContent strings.Builder
	var currentFileName string
	var currentDirectory string
	var parseError error
	currentFileOptions := make(map[string]string)
	symlinks = make(map[string]string)
	globalOptions = make(map[string]string)

	for _, line := range lines {
		ok := parseSymlinkFromTest(line, symlinks)
		if ok {
			continue
		}
		if testMetaData := optionRegex.FindStringSubmatch(line); testMetaData != nil {
			// Comment line, check for global/file @options and record them
			metaDataName := strings.ToLower(testMetaData[1])
			metaDataValue := strings.TrimSpace(testMetaData[2])
			if metaDataName == "currentdirectory" {
				currentDirectory = metaDataValue
			}
			if metaDataName != "filename" {
				if slices.Contains(fourslashDirectives, metaDataName) {
					// File-specific option
					currentFileOptions[metaDataName] = metaDataValue
				} else {
					// Global option
					if existingValue, ok := globalOptions[metaDataName]; ok && existingValue != metaDataValue {
						// !!! This would break existing submodule tests
						// panic("Duplicate global option: " + metaDataName)
					}
					globalOptions[metaDataName] = metaDataValue
				}
				continue
			}

			// New metadata statement after having collected some code to go with the previous metadata
			if currentFileName != "" {
				// Store result file
				newTestFile, e := parseFile(currentFileName, currentFileContent.String(), currentFileOptions)
				if e != nil {
					parseError = e
					break
				}
				testUnits = append(testUnits, newTestFile)

				// Reset local data
				currentFileContent.Reset()
				currentFileName = metaDataValue
				currentFileOptions = make(map[string]string)
			} else {
				// First metadata marker in the file
				currentFileName = strings.TrimSpace(testMetaData[2])
				if currentFileContent.Len() != 0 && scanner.SkipTrivia(currentFileContent.String(), 0) != currentFileContent.Len() {
					panic("Non-comment test content appears before the first '// @Filename' directive")
				}
				currentFileContent.Reset()
			}
		} else {
			// Subfile content line
			// Append to the current subfile content, inserting a newline if needed
			if currentFileContent.Len() != 0 {
				// End-of-line
				currentFileContent.WriteRune('\n')
			}
			currentFileContent.WriteString(line)
		}
	}

	// normalize the fileName for the single file case
	if len(testUnits) == 0 && len(currentFileName) == 0 {
		currentFileName = tspath.GetBaseFileName(fileName)
	}

	// if there are no parse errors so far, parse the rest of the file
	if parseError == nil {
		// EOF, push whatever remains
		newTestFile2, e := parseFile(currentFileName, currentFileContent.String(), currentFileOptions)

		parseError = e
		testUnits = append(testUnits, newTestFile2)
	}

	return testUnits, symlinks, currentDirectory, globalOptions, parseError
}

func extractCompilerSettings(content string) rawCompilerSettings {
	opts := make(map[string]string)

	for _, match := range optionRegex.FindAllStringSubmatch(content, -1) {
		opts[strings.ToLower(match[1])] = strings.TrimSuffix(strings.TrimSpace(match[2]), ";")
	}

	return opts
}

func parseSymlinkFromTest(line string, symlinks map[string]string) bool {
	linkMetaData := linkRegex.FindStringSubmatch(line)
	if len(linkMetaData) == 0 {
		return false
	}

	symlinks[strings.TrimSpace(linkMetaData[2])] = strings.TrimSpace(linkMetaData[1])
	return true
}
