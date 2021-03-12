package vsstoragecheck

import "testing"

func TestVSStorageCheck_CheckStorage_Capacity(t *testing.T) {
	fakeStorage := VSStorage{
		Id:            "test",
		State:         StorageStateReady,
		Type:          LocalStorage,
		Capacity:      10000,
		FreeCapacity:  5000,
		Timestamp:     "",
		Methods:       nil,
		LowWatermark:  0,
		HighWatermark: 8000,
	}
	c := VSStorageCheck{}

	results := c.CheckStorage(&fakeStorage, false)
	if len(results) != 0 {
		t.Errorf("Got unexpected alerts for storage at ok capacity: %v", results)
	}

	fakeStorage.FreeCapacity = 1000
	watermarkResults := c.CheckStorage(&fakeStorage, false)
	if len(watermarkResults) != 1 {
		t.Errorf("got unexpected alert count on over-watermark test, expected 1 got %d", len(watermarkResults))
	}

	fakeStorage.FreeCapacity = 2
	overCapResults := c.CheckStorage(&fakeStorage, false)
	if len(overCapResults) != 2 {
		t.Errorf("got unexpected alert count on over-capacity test, expected 2 got %d", len(overCapResults))
	}
}
