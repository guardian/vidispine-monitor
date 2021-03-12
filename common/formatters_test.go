package common

import "testing"

func TestFormatBytes(t *testing.T) {
	var kbValue int64 = 2048
	kbResult := FormatBytes(kbValue)
	if kbResult != "2Kib" {
		t.Errorf("2048 should return 2Kib, not %s", kbResult)
	}

	var mbValue int64 = 2048 * 1024
	mbResult := FormatBytes(mbValue)
	if mbResult != "2Mib" {
		t.Errorf("Expected 2Mib, got %s", mbResult)
	}

	var gbValue int64 = 2048 * 1024 * 1024
	gbResult := FormatBytes(gbValue)
	if gbResult != "2Gib" {
		t.Errorf("Expected 2Gib, got %s", gbResult)
	}
}
