package picamera

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/awfulbits/astro-raspicam/config"
	"github.com/dhowden/raspicam"
)

/***************************************************
Setting all the sensor configurations for each image
***************************************************/

// CameraConfiguration is exactly that
type CameraConfiguration struct {
	// This is a revised timeout per capture, or set, instead of per shot
	Timeout              int
	Sharpness            int
	Contrast             int
	Brightness           int
	Saturation           int
	ISO                  int
	ExposureCompensation int
	MeteringMode         raspicam.MeteringMode
	ShutterSpeed         int
	VideoStabilisation   bool
}

// ImageSetConfigSubset just wraps CameraConfiguration and Reps to be easily passed around
type ImageSetConfigSubset struct {
	Config CameraConfiguration
	Reps   int
}

// ImageSetConfig is a collection of subset configurations
type ImageSetConfig struct {
	ID      string
	Name    string
	Subsets []ImageSetConfigSubset
}

// Generate returns a ImageSet
func (imageSetConfig *ImageSetConfig) Generate() (set ImageSet) {
	var subsets []Subset
	for _, s := range imageSetConfig.Subsets {
		subset := Subset{
			raspicam.NewStill(),
			s.Reps,
		}

		// Set the configuration settings from subset
		subset.Still.Camera.Sharpness = s.Config.Sharpness
		subset.Still.Camera.Contrast = s.Config.Contrast
		subset.Still.Camera.Brightness = s.Config.Brightness
		subset.Still.Camera.Saturation = s.Config.Saturation
		subset.Still.Camera.ISO = s.Config.ISO
		subset.Still.Camera.ExposureCompensation = s.Config.ExposureCompensation
		subset.Still.Camera.MeteringMode = s.Config.MeteringMode
		subset.Still.Camera.ShutterSpeed = time.Duration(s.Config.ShutterSpeed) * time.Millisecond
		subset.Still.Camera.VideoStabilisation = s.Config.VideoStabilisation

		// Overwrite some default raspicam settings
		subset.Still.Quality = 100
		subset.Still.Raw = true
		subset.Still.Timeout = 100 * time.Millisecond

		// Add subset to the list of subsets
		subsets = append(subsets, subset)
	}

	set.Subsets = subsets
	set.ImageSetConfig = *imageSetConfig

	return
}

/*****************************************
Fully configured camera, ready for liftoff
*****************************************/

// This is not needed yet, but will be
// Camera is a, or a set of, sensor configuration(s) which impliment the Capture() method
// type Camera interface {
// 	Capture()
// }

// Subset is just a wrapper for camera configurations
type Subset struct {
	Still *raspicam.Still
	Reps  int
}

// ImageSet is the camera module, containing subsets of camera configurations
type ImageSet struct {
	Subsets []Subset
	ImageSetConfig
}

// Capture will use Set data to produce a set of images
func (set *ImageSet) Capture() (capErr error) {
	setPath := filepath.Join(config.C.ImageSetsPath, time.Now().Format("20060102150405"))
	if _, err := os.Stat(setPath); os.IsNotExist(err) {
		os.Mkdir(setPath, 0755)
	}

	for si, ss := range set.Subsets {
		subsetPath := filepath.Join(setPath, "sub_set-"+strconv.Itoa(si))
		os.Mkdir(subsetPath, 0755)

		for i := 0; i < ss.Reps; i++ {
			var file *os.File
			timestamp := time.Now().Format("20060102150405.000")
			file, err := os.Create(filepath.Join(subsetPath, timestamp+".jpg"))
			if err != nil {
				return
			}
			defer file.Close()

			errCh := make(chan error)
			go func() {
				for e := range errCh {
					capErr = e
				}
			}()

			raspicam.Capture(ss.Still, file, errCh)
			if capErr != nil {
				return
			}
		}
	}

	return
}
