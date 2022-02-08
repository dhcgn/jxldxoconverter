package jxlhandler

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhcgn/jxldxoconverter/config"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	//go:embed ..\assets\cjxl.exe
	cjxl_executable_data []byte

	compatible_exetentions = []string{"png", "apng", "gif", "jpeg", "jpg", "ppm", "pfm", "pgx"}
)

func IsCompatible(path string) bool {
	compatible := false
	for _, ext := range compatible_exetentions {
		if strings.HasSuffix(strings.ToLower(os.Args[1]), ext) {
			compatible = true
		}
	}
	return compatible
}

func ConvertToJxl(input, output, workingDir string, ffs config.FileFormatSetting, log *log.Entry) {
	cjxlTempPath := filepath.Join(workingDir, "cjxl.exe")
	fi, err := os.Stat(cjxlTempPath)
	if os.IsNotExist(err) || fi.Size() != int64(len(cjxl_executable_data)) {
		log.Println("Writing cjxl.exe")
		os.WriteFile(cjxlTempPath, cjxl_executable_data, 0644)
	}

	var cmd *exec.Cmd
	if ffs.DefaultConfig {
		log.Info("Using default config, no FileFormatSetting from config.json was found for this file")
		cmd = exec.Command(cjxlTempPath, input, output)
	} else {
		log.WithFields(logrus.Fields{"quality": ffs.Quality, "effort": ffs.Effort}).Info("Use FileFormatSetting from config.json")
		cmd = exec.Command(cjxlTempPath, input, output, "-q", fmt.Sprint(ffs.Quality), "-e", fmt.Sprint(ffs.Effort))
	}

	w := log.Writer()
	defer w.Close()
	cmd.Stderr = w
	cmd.Stdout = w

	start := time.Now()
	if err := cmd.Run(); err != nil {
		log.Println("Error: ", err)
	}
	log.Println("Convert to jxl in", time.Since(start))
}
