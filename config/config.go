package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// C is the config struct
var C = config{}

func init() {
	err := C.init()
	if err != nil {
		panic(err)
	}
}

type config struct {
	ImageSetConfigsPath string `toml:"image_set_configs_path"`
	ImageSetsPath       string `toml:"image_sets_path"`
}

func (c *config) init() error {
	_, err := toml.DecodeFile("config.toml", c)
	if err != nil {
		return err
	}

	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)

	if c.ImageSetConfigsPath == "" {
		c.ImageSetConfigsPath = filepath.Join(exPath, "image_set_configs")
	}
	if err := checkDir(c.ImageSetConfigsPath); err != nil {
		return err
	}

	if c.ImageSetsPath == "" {
		c.ImageSetsPath = filepath.Join(exPath, "image_sets")
	}
	if err := checkDir(c.ImageSetsPath); err != nil {
		return err
	}

	return nil
}

func checkDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
