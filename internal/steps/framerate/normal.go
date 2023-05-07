package framerate

import (
	"strconv"
	"strings"
)

func newNormalSelection() selector {
	return selector{
		Name:        "Normal",
		Description: "Select the framerate, e.g. <30> frames per second)",
		Prompt:      "Framerate (img/s)",
		ToFramerate: func(s string) (float64, error) {
			return strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
		},
	}
}
