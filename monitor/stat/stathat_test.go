package stat

import (
	"os"
	"testing"
)

func Test_NewStatHat(t *testing.T) {
	os.Setenv("SH_KEY", "helloKey")

	s := NewStatHat("foo")
	// it overrides the environment when non-empty key passed in
	if s.key != "foo" {
		t.Fatalf("Expected: %q, got: %v\n", "foo", s.key)
	}

	// it uses environment when empty key passed in
	s = NewStatHat("")
	if s.key != "helloKey" {
		t.Fatalf("Expected: %q, got: %v\n", "helloKey", s.key)
	}
}
