package slimdown

import (
	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

func testInit() {
	if err := globalFilepath.Init(); err != nil {
		panic(err)
	}
}
