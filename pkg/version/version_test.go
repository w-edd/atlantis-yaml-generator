package version

import "testing"

func TestGetVersion(t *testing.T) {
	expectedVersion := VERSION
	actualVersion := GetVersion()
	if actualVersion != expectedVersion {
		t.Errorf("Expected version %s, but got %s", expectedVersion, actualVersion)
	}
}
