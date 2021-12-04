package main

import "testing"

func TestRun_main(t *testing.T) {
	_, err := run_main()
	if err != nil {
		t.Error("failed run_main()")
	}
}
