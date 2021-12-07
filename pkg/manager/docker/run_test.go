package docker

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	got := Run("ubuntu", "latest", []string{"whoami"})
	gotStripped := strings.Trim(got, "\n ")
	want := "root"

	if gotStripped != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
