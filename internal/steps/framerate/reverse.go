package framerate

import (
	"strconv"
	"strings"
)

func newReverseSelection() selector {
	return selector{
		Name:        "Reverse",
		Description: "Select the duration an image is shown, e.g. <30>ms per image)",
		Prompt:      "Image screentime (ms/img)",
		ToFramerate: func(s string) (float64, error) {
			duration, err := strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
			if err != nil {
				return -1, err
			}
			return 1000 / duration, nil
		},
	}
}
