package apache

import (
  "testing"
)

func TestInit(t *testing.T) {
  inst := new(Apache22)
  inst.Init("foo", "bar")

  if inst.kind != "foo" && inst.kind != "bar" {
    t.Errorf("Expected foo and bar, got %s and %s\n", inst.kind, inst.mainUrl)
  }
}

func TestGetStatusForAll(t *testing.T) {

}
