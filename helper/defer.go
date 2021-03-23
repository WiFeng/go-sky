package helper

var (
	deferFunc []func()
)

// AddDeferFunc ...
func AddDeferFunc(f func()) {
	// no op
	if f == nil {
		return
	}

	deferFunc = append(deferFunc, f)
}

// RunDeferFunc ...
func RunDeferFunc() {
	for _, f := range deferFunc {
		f()
	}
}
