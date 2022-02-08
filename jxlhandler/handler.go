package jxlhandler

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

func ConvertToJxl(input, output, workingDir string) {
	cjxlTempPath := filepath.Join(workingDir, "cjxl.exe")
	fi, err := os.Stat(cjxlTempPath)
	if os.IsNotExist(err) || fi.Size() != int64(len(cjxl_executable_data)) {
		fmt.Println("Writing cjxl.exe")
		os.WriteFile(cjxlTempPath, cjxl_executable_data, 0644)
	}

	start := time.Now()
	cmd := exec.Command(cjxlTempPath, input, output)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Convert to jxl in", time.Since(start))
}
