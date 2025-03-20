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
		fmt.Printf("> %s [Y/n]: ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\t\n ") // I presume, at this moment in time, that this is enough

		if text == "Y" || text == "y" {
			return true
		} else if text == "N" || text == "n" {
			return false
		} else {
			// just ask again
		}
	}
}
