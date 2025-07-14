package diagnostic

import (
	"sync"
)

// Reporter collects and manages diagnostics.
type Reporter struct {
	mu          sync.Mutex
	diagnostics []*Diagnostic
}

func NewReporter() *Reporter {
	return &Reporter{
		diagnostics: make([]*Diagnostic, 0),
	}
}

func (r *Reporter) Add(d *Diagnostic) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.diagnostics = append(r.diagnostics, d)
}

func (r *Reporter) All() []*Diagnostic {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.diagnostics
}

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

func (r *Reporter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.diagnostics = make([]*Diagnostic, 0)
}
