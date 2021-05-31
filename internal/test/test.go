package test

import (
	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

func init() {
	if err := globalFilepath.Init("assets"); err != nil {
		panic(err)
	}
}
