// Package config contains models for user settings
package config

import (
	"os"
	"path/filepath"
	"time"
)

type config struct {
	TimeFrom         time.Time
	TimeTo           time.Time
	AbsInputDir      string
	AbsOutpuFilepath string
	ImageNamesSorted []string
	ImgsPerSecond    float64
}

// Cfg is the current app configuration
var Cfg = config{}

// TODO: change to Linux once ready for Docker
const (
	ImagesInputRootDir string = "D:/timelapse/input"
	// ImagesInputRootDir  string = "Y:/"
	VideosOutputRootDir string = "D:/timelapse/output"
)

// Init prepares the required in- and output directories.
func Init() (err error) {
	if err = os.MkdirAll(ImagesInputRootDir, 0644); err != nil {
		return
	}
	err = os.MkdirAll(VideosOutputRootDir, 0644)
	return
}

// SetInputDir saves the input directory to the current config.
// No validation is performed.
func SetInputDir(d string) {
	Cfg.AbsInputDir = filepath.Join(ImagesInputRootDir, d)
}

// SetOutputFilename saves the output filename to the current config.
// No validation is performed.
func SetOutputFilename(f string) {
	Cfg.AbsOutpuFilepath = filepath.Join(VideosOutputRootDir, f)
}

// SetTimerange saves the input photo timerange to the current config.
// No validation is performed.
func SetTimerange(from time.Time, to time.Time) {
	Cfg.TimeFrom = from
	Cfg.TimeTo = to
}

// SetImagesPerSecond saves the framerate to the current config.
// No validation is performed.
func SetImagesPerSecond(framerate float64) {
	Cfg.ImgsPerSecond = framerate
}

// GetApproxTotalDuration uses the configured framerate and list of image filenames to predict
// the approximate final timelapse video duration.
func GetApproxTotalDuration() time.Duration {
	numFiles := len(Cfg.ImageNamesSorted)
	imgsPerSecond := Cfg.ImgsPerSecond
	return time.Second * time.Duration(float64(numFiles)/imgsPerSecond)
}
