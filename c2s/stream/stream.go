package stream

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/kpmy/xippo/tools/srv"
	"github.com/kpmy/xippo/units"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/halt"
	"hash/adler32"
	"net"
	"reflect"
	"time"
)

type Stream interface {
	Server() *units.Server
	Write(*bytes.Buffer) error
	Ring(func(*bytes.Buffer) bool, time.Duration)
	Id(...string) string
	Close()
}

type wrapperStream struct {
	base        Stream
	onErrorFunc func(error)
}

func (w *wrapperStream) Write(b *bytes.Buffer) error { return w.base.Write(b) }

func (w *wrapperStream) Ring(fn func(*bytes.Buffer) bool, t time.Duration) {
	w.base.Ring(fn, t)
}

func (w *wrapperStream) Server() *units.Server { return w.base.Server() }

func (w *wrapperStream) Id(s ...string) string { return w.base.Id(s...) }

func (w *wrapperStream) Close() { w.base.Close() }

type dummyStream struct {
	to *units.Server
}

func (d *dummyStream) Ring(func(*bytes.Buffer) bool, time.Duration) { panic(126) }
func (d *dummyStream) Write(b *bytes.Buffer) error                  { panic(126) }
func (d *dummyStream) Server() *units.Server                        { return d.to }
func (d *dummyStream) Id(...string) string                          { return "" }
func (d *dummyStream) Close()                                       {}

type xmppStream struct {
	to   *units.Server
	conn net.Conn
	ctrl chan bool
	data chan pack

	onCloseFunc func(error)

	id  string
	jid string
}

type pack struct {
	data []byte
	hash uint32
}

func (x *xmppStream) Id(s ...string) string {
	if len(s) > 0 {
		x.id = s[0]
	}
	return x.id
}
func (x *xmppStream) Server() *units.Server { return x.to }

func (x *xmppStream) Write(b *bytes.Buffer) (err error) {
	//log.Println("OUT")
	//log.Println(string(b.Bytes()))
	//log.Println()
	_, err = x.conn.Write(b.Bytes())
	return
}

func (x *xmppStream) Ring(fn func(*bytes.Buffer) bool, timeout time.Duration) {
	timed := make(chan bool)
	if timeout > 0 {
		go func() {
			<-time.NewTimer(timeout).C
			timed <- true
		}()
	}
	for stop := false; !stop; {
		select {
		case p := <-x.data:
			done := fn(bytes.NewBuffer(p.data))
			if !done {
				x.data <- pack{data: p.data, hash: p.hash}
			} else {
				stop = true
			}
		case <-timed:
			stop = true
		}
	}
}

func (x *xmppStream) closeStream() {
	assert.For(x.ctrl != nil, 20)
	x.ctrl <- true
}

func (x *xmppStream) onClose(err error) {
	if x.onCloseFunc != nil {
		x.onCloseFunc(err)
	}
}

func (x *xmppStream) Close() {
	x.closeStream()
	x.onClose(nil)
}

func New(to *units.Server, onErr ...func(error)) Stream {
	var onErrFunc func(error)
	if len(onErr) > 0 {
		onErrFunc = onErr[0]
	}
	return &wrapperStream{base: &dummyStream{to: to}, onErrorFunc: onErrFunc}
}

func Bind(_s Stream, jid string) {
	switch w := _s.(type) {
	case *wrapperStream:
		switch s := w.base.(type) {
		case *xmppStream:
			s.jid = jid
		}
	}
}

func Dial(_s Stream) (err error) {
	switch w := _s.(type) {
	case *wrapperStream:
		switch s := w.base.(type) {
		case *dummyStream:
			x := &xmppStream{to: s.to, onCloseFunc: w.onErrorFunc}
			var (
				host, port string
				dial       net.Dialer
			)
			if host, port, err = srv.Resolve(x.to); err == nil {
				dial.KeepAlive = 10 * time.Second
				if x.conn, err = dial.Dial("tcp", host+":"+port); err == nil {
					x.ctrl = make(chan bool)
					x.data = make(chan pack, 256)
					go func(stream *xmppStream) {
						<-stream.ctrl
						stream.conn.Close()
					}(x)

					doSplit := func(stream *xmppStream, onErr func(error)) {
						for data := range spl1t(stream.conn, onErr) {
							//log.Println("SPLIT")
							//log.Println(string(data))
							//log.Println()
							stream.data <- pack{data: data, hash: adler32.Checksum(data)}
						}
					}

					var onErr func(error)
					onErr = func(_err error) {
						switch e := _err.(type) {
						case net.Error:
							if e.Temporary() {
								go doSplit(x, onErr)
							} else {
								x.closeStream()
								x.onClose(e)
								w.base = s
							}
						default:
							panic(fmt.Sprint(reflect.TypeOf(e), e))
						}
					}

					go doSplit(x, onErr)
					w.base = x
				} else {
					x.onClose(err)
				}
			} else {
				x.onClose(err)
			}
		default:
			err = errors.New("already connected")
		}
	default:
		halt.As(100, reflect.TypeOf(_s))
	}
	return
}
