package diagnostic

import (
	"sync"
)

// Reporter collects and manages diagnostics.
type Reporter struct {
	mu          sync.Mutex
	diagnostics []*Diagnostic
}

// NewReporter creates a new reporter.
func NewReporter() *Reporter {
	return &Reporter{
		diagnostics: make([]*Diagnostic, 0),
	}
}

// Add adds a new diagnostic to the reporter.
func (r *Reporter) Add(d *Diagnostic) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.diagnostics = append(r.diagnostics, d)
}

// All returns all diagnostics.
func (r *Reporter) All() []*Diagnostic {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.diagnostics
}

// Errors returns all error diagnostics.
func (r *Reporter) Errors() []*Diagnostic {
	r.mu.Lock()
	defer r.mu.Unlock()
	var errors []*Diagnostic
	for _, d := range r.diagnostics {
		if d.Severity == SeverityError {
			errors = append(errors, d)
		}
	}
	return errors
}

// Warnings returns all warning diagnostics.
func (r *Reporter) Warnings() []*Diagnostic {
	r.mu.Lock()
	defer r.mu.Unlock()
	var warnings []*Diagnostic
	for _, d := range r.diagnostics {
		if d.Severity == SeverityWarning {
			warnings = append(warnings, d)
		}
	}
	return warnings
}

// Infos returns all informational diagnostics.
func (r *Reporter) Infos() []*Diagnostic {
	r.mu.Lock()
	defer r.mu.Unlock()
	var infos []*Diagnostic
	for _, d := range r.diagnostics {
		if d.Severity == SeverityInfo {
			infos = append(infos, d)
		}
	}
	return infos
}

// HasErrors returns true if there are any errors.
func (r *Reporter) HasErrors() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, d := range r.diagnostics {
		if d.Severity == SeverityError {
			return true
		}
	}
	return false
}

// HasWarnings returns true if there are any warnings.
func (r *Reporter) HasWarnings() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, d := range r.diagnostics {
		if d.Severity == SeverityWarning {
			return true
		}
	}
	return false
}

// Reset clears all diagnostics from the reporter.
func (r *Reporter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.diagnostics = make([]*Diagnostic, 0)
}
