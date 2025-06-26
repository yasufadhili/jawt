package compiler

import "fmt"

type Compiler struct {
	FileType string
}

func NewCompiler(fileType string) (*Compiler, error) {
	if fileType != "Component" && fileType != "Page" {
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}
	return &Compiler{
		FileType: fileType,
	}, nil
}

func (c *Compiler) Compile() (string, error) {

	return "", nil
}
