package dex

import "fmt"

func generateURL(apex string) string {
	return fmt.Sprintf("dex.%s", apex)
}
