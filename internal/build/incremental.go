package build

import "fmt"

// BuildIncremental performs incremental build with error handling
func (b *BuildSystem) BuildIncremental() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.compilerManager == nil {
		return fmt.Errorf("compiler not initialised - run full build first")
	}

	if err := b.compilerManager.CompileChanged(); err != nil {
		buildErr := fmt.Errorf("incremental compilation failed: %w", err)
		if b.errorState.shouldShowError(buildErr) {
			b.printError("Incremental Compilation", buildErr)
		}
		return buildErr
	}

	if b.errorState.shouldShowError(nil) {
		b.printIncrementalSuccess()
	}

	return nil
}
