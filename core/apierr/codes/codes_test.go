package codes

import "testing"

func TestCodeString(t *testing.T) {
	for want := range nCodes {
		if got := FromString(want.String()); got != want {
			t.Errorf("FromString(%q) = %q, want %q", want.String(), got, want)
		}
	}
}
