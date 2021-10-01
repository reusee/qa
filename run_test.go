package qa

import "testing"

func TestRun(t *testing.T) {
	errs := Run(".")
	if len(errs) > 0 {
		t.Fatal()
	}
}
