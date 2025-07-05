package diagnostic

import (
	"sync"
)

// Reporter collects and manages diagnostics.
type Reporter struct {
	mu       sync.Mutex
	errors   []*Error
	warnings []*Error
	infos    []*Error
}

// NewReporter creates a new reporter.
func NewReporter() *Reporter {
	return &Reporter{
		errors:   make([]*Error, 0),
		warnings: make([]*Error, 0),
		infos:    make([]*Error, 0),
	}
}

// Add adds a new diagnostic to the reporter.
func (r *Reporter) Add(err *Error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch err.Severity {
	case SeverityError:
		r.errors = append(r.errors, err)
	case SeverityWarning:
		r.warnings = append(r.warnings, err)
	case SeverityInfo:
		r.infos = append(r.infos, err)
	}
}

// Errors returns all errors.
func (r *Reporter) Errors() []*Error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.errors
}

// Warnings returns all warnings.
func (r *Reporter) Warnings() []*Error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.warnings
}

// Infos returns all informational messages.
func (r *Reporter) Infos() []*Error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.infos
}

// HasErrors returns true if there are any errors.
func (r *Reporter) HasErrors() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.errors) > 0
}

// HasWarnings returns true if there are any warnings.
func (r *Reporter) HasWarnings() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.warnings) > 0
}

// Reset clears all diagnostics from the reporter.
func (r *Reporter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.errors = make([]*Error, 0)
	r.warnings = make([]*Error, 0)
	r.infos = make([]*Error, 0)
}
