package assets

// SetTestRootDir overrides the root directory for tests.
// It returns a restore function.
func SetTestRootDir(dir string) func() {
	old := defaultDirs.rootDir
	defaultDirs.rootDir = dir

	return func() {
		defaultDirs.rootDir = old
	}
}
