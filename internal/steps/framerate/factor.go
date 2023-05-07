package framerate

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stnokott/go-timelapse/internal/config"
)

func newFactorSelection() selector {
	return selector{
		Name:        "Factor",
		Description: "Select a speedup factor, relative to the time from first to last configured image.",
		Prompt:      "Speedup factor",
		ToFramerate: func(s string) (f float64, err error) {
			var factor float64
			factor, err = strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
			if err != nil {
				return
			}
			numFiles := len(config.Cfg.ImageNamesSorted)
			firstImg := filepath.Join(config.Cfg.AbsInputDir, config.Cfg.ImageNamesSorted[0])
			lastImg := filepath.Join(config.Cfg.AbsInputDir, config.Cfg.ImageNamesSorted[numFiles-1])
			var fiFirst, fiLast fs.FileInfo
			fiFirst, err = os.Stat(firstImg)
			if err != nil {
				return
			}
			fiLast, err = os.Stat(lastImg)
			if err != nil {
				return
			}
			duration := fiLast.ModTime().Sub(fiFirst.ModTime())
			f = float64(numFiles) / (duration.Seconds() / factor)
			return
		},
	}
}
