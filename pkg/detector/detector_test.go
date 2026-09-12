package detector

import (
	"testing"
)

func TestDetect(t *testing.T) {
	ctx := Detect()

	if ctx.OS == "" {
		t.Errorf("Expected OS to be populated, got empty string")
	}
	if ctx.Arch == "" {
		t.Errorf("Expected Arch to be populated, got empty string")
	}
	if len(ctx.ActiveTargets) == 0 {
		t.Errorf("Expected at least one active target, got 0")
	}
}
