package helper

var (
	deferFunc []func()
)

// AddDeferFunc ...
func AddDeferFunc(f func()) {
	deferFunc = append(deferFunc, f)
}

// RunDeferFunc ...
func RunDeferFunc() {
	for _, f := range deferFunc {
		f()
	}
}
