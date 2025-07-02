package project

import (
	"context"
	"fmt"
	"iter"
	"sync"

	"github.com/yasufadhili/jawt/internal/tsc/ast"
	"github.com/yasufadhili/jawt/internal/tsc/checker"
	"github.com/yasufadhili/jawt/internal/tsc/compiler"
	"github.com/yasufadhili/jawt/internal/tsc/core"
)

type checkerPool struct {
	maxCheckers int
	program     *compiler.Program

	mu                  sync.Mutex
	cond                *sync.Cond
	createCheckersOnce  sync.Once
	checkers            []*checker.Checker
	inUse               map[*checker.Checker]bool
	fileAssociations    map[*ast.SourceFile]int
	requestAssociations map[string]int
	log                 func(msg string)
}

var _ compiler.CheckerPool = (*checkerPool)(nil)

func newCheckerPool(maxCheckers int, program *compiler.Program, log func(msg string)) *checkerPool {
	pool := &checkerPool{
		program:             program,
		maxCheckers:         maxCheckers,
		checkers:            make([]*checker.Checker, maxCheckers),
		inUse:               make(map[*checker.Checker]bool),
		requestAssociations: make(map[string]int),
		log:                 log,
	}

	pool.cond = sync.NewCond(&pool.mu)
	return pool
}

func (p *checkerPool) GetCheckerForFile(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func()) {
	p.mu.Lock()
	defer p.mu.Unlock()

	requestID := core.GetRequestID(ctx)
	if requestID != "" {
		if checker, release := p.getRequestCheckerLocked(requestID); checker != nil {
			return checker, release
		}
	}

	if p.fileAssociations == nil {
		p.fileAssociations = make(map[*ast.SourceFile]int)
	}

	if index, ok := p.fileAssociations[file]; ok {
		checker := p.checkers[index]
		if checker != nil {
			if inUse := p.inUse[checker]; !inUse {
				p.inUse[checker] = true
				if requestID != "" {
					p.requestAssociations[requestID] = index
				}
				return checker, p.createRelease(requestID, index, checker)
			}
		}
	}

	checker, index := p.getCheckerLocked(requestID)
	p.fileAssociations[file] = index
	return checker, p.createRelease(requestID, index, checker)
}

func (p *checkerPool) GetChecker(ctx context.Context) (*checker.Checker, func()) {
	p.mu.Lock()
	defer p.mu.Unlock()
	checker, index := p.getCheckerLocked(core.GetRequestID(ctx))
	return checker, p.createRelease(core.GetRequestID(ctx), index, checker)
}

func (p *checkerPool) Files(checker *checker.Checker) iter.Seq[*ast.SourceFile] {
	panic("unimplemented")
}

func (p *checkerPool) GetAllCheckers(ctx context.Context) ([]*checker.Checker, func()) {
	p.mu.Lock()
	defer p.mu.Unlock()

	requestID := core.GetRequestID(ctx)
	if requestID == "" {
		panic("cannot call GetAllCheckers on a project.checkerPool without a request ID")
	}

	// A request can only access one checker
	if c, release := p.getRequestCheckerLocked(requestID); c != nil {
		return []*checker.Checker{c}, release
	}

	c, release := p.GetChecker(ctx)
	return []*checker.Checker{c}, release
}

func (p *checkerPool) getCheckerLocked(requestID string) (*checker.Checker, int) {
	if checker, index := p.getImmediatelyAvailableChecker(); checker != nil {
		p.inUse[checker] = true
		if requestID != "" {
			p.requestAssociations[requestID] = index
		}
		return checker, index
	}

	if !p.isFullLocked() {
		checker, index := p.createCheckerLocked()
		p.inUse[checker] = true
		if requestID != "" {
			p.requestAssociations[requestID] = index
		}
		return checker, index
	}

	checker, index := p.waitForAvailableChecker()
	p.inUse[checker] = true
	if requestID != "" {
		p.requestAssociations[requestID] = index
	}
	return checker, index
}

func (p *checkerPool) getRequestCheckerLocked(requestID string) (*checker.Checker, func()) {
	if index, ok := p.requestAssociations[requestID]; ok {
		checker := p.checkers[index]
		if checker != nil {
			if inUse := p.inUse[checker]; !inUse {
				p.inUse[checker] = true
				return checker, p.createRelease(requestID, index, checker)
			}
			// Checker is in use, but by the same request - assume it's the
			// same goroutine or is managing its own synchronization
			return checker, noop
		}
	}
	return nil, noop
}

func (p *checkerPool) getImmediatelyAvailableChecker() (*checker.Checker, int) {
	for i, checker := range p.checkers {
		if checker == nil {
			continue
		}
		if inUse := p.inUse[checker]; !inUse {
			return checker, i
		}
	}

	return nil, -1
}

func (p *checkerPool) waitForAvailableChecker() (*checker.Checker, int) {
	p.log("checkerpool: Waiting for an available checker")
	for {
		p.cond.Wait()
		checker, index := p.getImmediatelyAvailableChecker()
		if checker != nil {
			return checker, index
		}
	}
}

func (p *checkerPool) createRelease(requestId string, index int, checker *checker.Checker) func() {
	return func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		delete(p.requestAssociations, requestId)
		if checker.WasCanceled() {
			// Canceled checkers must be disposed
			p.log(fmt.Sprintf("checkerpool: Checker for request %s was canceled, disposing it", requestId))
			p.checkers[index] = nil
			delete(p.inUse, checker)
		} else {
			p.inUse[checker] = false
		}
		p.cond.Signal()
	}
}

func (p *checkerPool) isFullLocked() bool {
	for _, checker := range p.checkers {
		if checker == nil {
			return false
		}
	}
	return true
}

func (p *checkerPool) createCheckerLocked() (*checker.Checker, int) {
	for i, existing := range p.checkers {
		if existing == nil {
			checker := checker.NewChecker(p.program)
			p.checkers[i] = checker
			return checker, i
		}
	}
	panic("called createCheckerLocked when pool is full")
}

func (p *checkerPool) isRequestCheckerInUse(requestID string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if index, ok := p.requestAssociations[requestID]; ok {
		checker := p.checkers[index]
		if checker != nil {
			return p.inUse[checker]
		}
	}
	return false
}

func (p *checkerPool) size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	size := 0
	for _, checker := range p.checkers {
		if checker != nil {
			size++
		}
	}
	return size
}

func noop() {}
