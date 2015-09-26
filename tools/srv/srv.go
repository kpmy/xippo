package srv

import (
	"xippo/units"
)

func Resolve(s *units.Server) (host, port string, err error) {
	return s.Name, "5222", nil
}
