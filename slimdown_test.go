package slimdown

import (
	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

func init() {
	if err := globalFilepath.Init(); err != nil {
		panic(err)
	}
}
