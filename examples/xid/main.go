package main

import (
	"github.com/rs/xid"
)

func main() {
	guid := xid.New()
	println(guid.Machine())
	println(guid.Pid())
	println(guid.Time().String())
	println(guid.Counter())
	println(guid.String())
}
