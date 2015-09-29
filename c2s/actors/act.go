package actors

import (
	"errors"
	"github.com/kpmy/xippo/c2s/stream"
	"github.com/kpmy/ypk/act"
)

func With() act.Continue {
	return act.Seq()
}
func C(fn func(stream.Stream) error) func(interface{}) (interface{}, error) {
	return func(x interface{}) (ret interface{}, err error) {
		if s, ok := x.(stream.Stream); ok {
			ret = s
			err = fn(s)
		} else {
			err = errors.New("unknown stream")
		}
		return
	}
}
