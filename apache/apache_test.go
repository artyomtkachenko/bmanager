package apache

import (
	"testing"
)

func TestNew(t *testing.T) {
	var inst Apache22
	inst.New("apache", "/balancer-manager")

	if inst.kind != "apache" {
		t.Errorf("Expected apache, got %s\n", inst.kind)
	}
	if inst.mainUrl != "/balancer-manager" {
		t.Errorf("Expected /balancer-manager, got %s\n", inst.mainUrl)
	}
}

func TestgetDetailsFromUri(t *testing.T) {
}
