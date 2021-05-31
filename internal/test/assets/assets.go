package assets

import (
	"os"

	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

func init() {
	if err := globalFilepath.Init("storage"); err != nil {
		panic(err)
	}
}

func Load(fileName string) (content []byte) {
	var err error

	content, err = os.ReadFile(globalFilepath.Join(fileName))
	if err != nil {
		panic(err)
	}

	return
}
