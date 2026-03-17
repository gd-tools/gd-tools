package model

import "testing"

func TestMountListHasMount(t *testing.T) {
	mounts := MountList{
		{
			Provider: "hetzner",
			ID:       "123",
			Path:     "/data",
		},
	}

	if !mounts.HasMount("/data") {
		t.Fatal("expected mount /data")
	}

	if mounts.HasMount("/other") {
		t.Fatal("unexpected mount /other")
	}
}
