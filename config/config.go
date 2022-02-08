package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/dhcgn/jxldxoconverter/helper"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	FileFormatSettings []FileFormatSetting `json:"file_formats"`
}

type FileFormatSetting struct {
	Extension        string `json:"extension"`
	Quality          int    `json:"quality"`
	Effort           int    `json:"effort"`
	DeleteSourceFile bool   `json:"delete_source_file"`
	Comment          string `json:"comment"`
	DefaultConfig    bool   `json:"default_config"`
}

func (c Config) GetFileFormatSetting(path string) FileFormatSetting {
	ext := filepath.Ext(path)

	for _, ffs := range c.FileFormatSettings {
		m, err := regexp.MatchString(ffs.Extension, ext)
		if err != nil {
			log.Error(err)
		}
		if m {
			return ffs
		}
	}

	return FileFormatSetting{
		DefaultConfig: true,
	}
}

func GetConfig(rootDir string) Config {
	configPath := filepath.Join(rootDir, "config.json")
	if !helper.Exists(configPath) {
		j, _ := json.MarshalIndent(Config{
			FileFormatSettings: []FileFormatSetting{
				{
					Extension:        "tif|tiff",
					Quality:          99,
					Effort:           8,
					DeleteSourceFile: true,
					Comment:          "tif files are created from dxo for this export, so they can be deleted afterwards. Effort 99 is best quality after loseless.",
				},
				{
					Extension:     "jpg|jpeg",
					DefaultConfig: true,
					Comment:       "Use defaults of JPEG XL encoder, JPGs will be converted to JXL LOSSLESS. No generation loss!",
				},
			},
		}, "", "  ")
		err := ioutil.WriteFile(configPath, j, 0644)
		if err != nil {
			log.Error(err)
		}
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Error(err)
	}

	var c Config
	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Error(err)
	}
	return c
}
