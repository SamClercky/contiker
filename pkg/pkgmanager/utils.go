package pkgmanager

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ask the question until a suitable answer has been given
func askPermission(question string) bool {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("> %s [Y/n]: ", question)
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\t\n ") // I presume, at this moment in time, that this is enough

		switch text {
		case "Y", "y":
			return true
		case "N", "n":
			return false
		default:
			// just ask again
		}
	}
}
