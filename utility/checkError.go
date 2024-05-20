package utility

import (
	"fmt"
	"os"
)

func CheckError(e error, customMessage string) {
	if e != nil {
		fmt.Printf("%s  %v", customMessage, e)
	}
	os.Exit(1)
}
