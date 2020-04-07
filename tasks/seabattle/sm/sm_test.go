package sm

import "testing"

func Test_Basics(t *testing.T) {
	machine := New()
	if machine.Current() != StateInit {
		t.Fatalf("incorrect state: want %s, got %s", StateInit, machine.Current())
	}
}
