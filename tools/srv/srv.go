package srv

import (
	"github.com/kpmy/xippo/units"
	"math"
	"net"
	"strconv"
)

const XMPP_CLIENT = "xmpp-client"

func Resolve(s *units.Server) (host, port string, erro error) {
	if _, addrs, err := net.LookupSRV(XMPP_CLIENT, "tcp", s.Name); err == nil {
		priority := uint16(math.MaxUint16)
		for _, s := range addrs {
			if s.Priority <= priority {
				host = s.Target
				port = strconv.Itoa(int(s.Port))
				priority = s.Priority
			}
		}
	} else {
		erro = err
	}
	return
}
