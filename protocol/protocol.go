package protocol

import (
	"fmt"
)

const DEBUG = true

func debug(v ...interface{}) {
	if DEBUG {
		fmt.Println(v...)
	}
}
