package dyn

import (
	"github.com/kpmy/xippo/entity"
	"github.com/kpmy/ypk/dom"
	"math/rand"
	"strconv"
)

type Entity interface {
	entity.Entity
	Type() string
}

const (
	PRESENCE = "presence"
	IQ       = "iq"
	MESSAGE  = "message"
	TYPE     = "type"
	ID       = "id"
	BODY     = "body"
	TO       = "to"
)

func NewMessage(typ entity.MessageType, to string, body string) (ret dom.Element) {
	ret = dom.Elem(MESSAGE)
	ret.Attr(TYPE, string(typ))
	ret.Attr(ID, strconv.FormatInt(int64(rand.Intn(0xffffff)), 16))
	ret.Attr(TO, to)
	b := dom.Elem(BODY)
	b.AppendChild(dom.Txt(body))
	ret.AppendChild(b)
	return
}
