package steps

import (
	"github.com/kpmy/xippo/c2s/stream"
	"github.com/kpmy/xippo/entity"
)

func InitialPresence(s stream.Stream) (err error) {
	err = s.Write(entity.ProduceStatic(&entity.PresencePrototype))
	return
}

func PresenceTo(jid string, show entity.PresenceShow, status string) func(stream.Stream) error {
	return func(s stream.Stream) error {
		pr := &entity.Presence{}
		*pr = entity.PresencePrototype
		pr.To = jid
		pr.Show = show
		pr.Status = status
		return s.Write(entity.ProduceStatic(pr))
	}
}
