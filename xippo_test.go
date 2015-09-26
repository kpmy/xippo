package xippo

import (
	"testing"
	"xippo/c2s/stream"
	"xippo/units"
)

func TestNothing(t *testing.T) {
	s := &units.Server{Name: "test"}
	stream.New(s)
}
