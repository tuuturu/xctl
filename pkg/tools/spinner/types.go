package spinner

import (
	"io"
	"time"

	"github.com/briandowns/spinner"
)

// Interesting charsets: 7 , 13 , 16 , 47 , 50 , 59 , 69 , 78
const (
	charsetIndex = 7
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
