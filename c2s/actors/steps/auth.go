package steps

import (
	"bytes"
	"github.com/kpmy/xippo/c2s/stream"
	"github.com/kpmy/xippo/entity"
	"github.com/kpmy/xippo/units"
	"log"
	"reflect"
)

type PlainAuth struct {
	Client *units.Client
	Pwd    string
}

func (p *PlainAuth) Act() func(stream.Stream) error {
	return func(s stream.Stream) (err error) {
		auth := &entity.PlainAuth{}
		*auth = entity.PlainAuthPrototype
		auth.Init(p.Client.Name, p.Pwd)
		if err = s.Write(entity.ProduceStatic(auth)); err == nil {
			s.Ring(func(b *bytes.Buffer) (done bool) {
				var _e entity.Entity
				if _e, err = entity.ConsumeStatic(b); err == nil {
					switch e := _e.(type) {
					case *entity.Success:
						done = true
					default:
						log.Println(reflect.TypeOf(e))
					}
				}
				return
			}, 0)
		}
		return
	}
}
