package spinner

import (
	"io"
	"time"

	"github.com/briandowns/spinner"
)

// Interesting charsets: 7 , 13 , 16 , 59 , 69
// 7: Deathstar 5/10
// 13: clean box 6/10
// 16: evil progress bar 7/10
// 59: dots 6/10
// 69: thicker dots 7/10
const (
	charsetIndex = 16
	speed        = 100 * time.Millisecond
	color        = "red"
)

func NewSpinner(out io.Writer) *spinner.Spinner {
	spin := spinner.New(spinner.CharSets[charsetIndex], speed,
		spinner.WithColor(color),
		spinner.WithWriter(out),
	)

	return spin
}
