package compiler

type AstBuilder struct {
}

func NewAstBuilder() *AstBuilder {
	return &AstBuilder{}
}

func (b *AstBuilder) BuildAST() (*JMLDocumentNode, error) {
	return nil, nil
}
