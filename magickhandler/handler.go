package magickhandler

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

//go:embed ..\assets\magick.exe
var magick_executable_data []byte

func ConvertToPng(source, workingDir string, log *log.Entry) string {
	magickTempPath := filepath.Join(workingDir, "magick.exe")
	fi, err := os.Stat(magickTempPath)
	if os.IsNotExist(err) || fi.Size() != int64(len(magick_executable_data)) {
		log.Println("Writing magick_temp.exe")
		os.WriteFile(magickTempPath, magick_executable_data, 0644)
	}

	ext := ".png"
	newFile := "temp_" + uuid.New().String() + ext
	newFile = filepath.Join(workingDir, newFile)
	start := time.Now()
	cmd := exec.Command(magickTempPath, "convert", source, newFile)
	w := log.Writer()
	defer w.Close()
	cmd.Stderr = w
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		log.Println("Error: ", err)
	}
	log.Println("Convert to png in", time.Since(start), "to", newFile)

	files, err := filepath.Glob(strings.TrimSuffix(newFile, ext) + "*")
	if err != nil {
		log.Println("Error: ", err)
	}

	for _, f := range files[1:] {
		os.Remove(f)
	}

	return files[0]
}
