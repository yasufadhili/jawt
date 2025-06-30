package common

// BuildTarget represents different output formats
type BuildTarget int

const (
	TargetPage BuildTarget = iota
	TargetComponent
)
