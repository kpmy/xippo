package srv

import (
	"github.com/kpmy/xippo/units"
	"testing"
)

func TestLookup(t *testing.T) {
	if host, port, err := Resolve(&units.Server{Name: "jabber.ru"}); err == nil {
		t.Log(host, port)
	} else {
		t.Error(err)
	}

	if host, port, err := Resolve(&units.Server{Name: "jabber.org"}); err == nil {
		t.Log(host, port)
	} else {
		t.Error(err)
	}
}
