package slimdown

import (
	"sync"

	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

var (
	testInitOnce sync.Once
)

func testInit() {
	testInitOnce.Do(func() {
		if err := globalFilepath.Init(); err != nil {
			panic(err)
		}
	})
}
