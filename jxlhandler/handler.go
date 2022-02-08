package jxlhandler

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

func ConvertToJxl(input, output, workingDir string, log *log.Entry) {
	cjxlTempPath := filepath.Join(workingDir, "cjxl.exe")
	fi, err := os.Stat(cjxlTempPath)
	if os.IsNotExist(err) || fi.Size() != int64(len(cjxl_executable_data)) {
		log.Println("Writing cjxl.exe")
		os.WriteFile(cjxlTempPath, cjxl_executable_data, 0644)
	}

	start := time.Now()
	cmd := exec.Command(cjxlTempPath, input, output)
	w := log.Writer()
	defer w.Close()
	cmd.Stderr = w
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		log.Println("Error: ", err)
	}
	log.Println("Convert to jxl in", time.Since(start))
}
