package log

import (
	"net"
	"time"
)

const (
	Error   string = "ERROR"
	Success string = "Success"
)

type Logger struct {
	Datetime time.Time
	IP       net.IP
	Params   string
	Status   string
}
