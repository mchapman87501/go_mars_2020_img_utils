package main

import "testing"

func TestIDsMatch(t *testing.T) {
	got := idsMatch("", "")
	if got {
		t.Errorf("Empty IDs should not match.")
	}

	got = idsMatch("12", "34")
	if got {
		t.Errorf("Short IDs should not match.")
	}

	got = idsMatch("NRF_1blah", "NLF_1blah")
	if !got {
		t.Error("IDs w. same prefix, and 'L' or 'R', should match.")
	}
}

