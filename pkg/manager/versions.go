package manager

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteCurrentVersion(tag, branch string) {
	envToWrite := []byte(fmt.Sprintf("CURRENT_TAG=%s\nCURRENT_BRANCH=%s", tag, branch))
	err := os.WriteFile(".env", envToWrite, 0666)
	check(err)
}
