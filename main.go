package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	//go:embed assets\cjxl.exe
	cjxl_executable_data   []byte
	compatible_exetentions = []string{"png", "apng", "gif", "jpeg", "jpg", "ppm", "pfm", "pgx"}

	//go:embed assets\magick.exe
	magick_executable_data []byte
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Need at least one Argument with a full path to a .tif file")
		os.Exit(1)
	}

	compatible := false
	for _, ext := range compatible_exetentions {
		if strings.HasSuffix(strings.ToLower(os.Args[1]), ext) {
			compatible = true
		}
	}

	rootDir := filepath.Dir(os.Args[0])
	workingDir := filepath.Join(rootDir, "temp")
	sourceFiles := os.Args[1:]

	if !exists(workingDir) {
		err := os.Mkdir(workingDir, 0755)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	for _, sourceFile := range sourceFiles {

		fi, err := os.Stat(sourceFile)
		if os.IsNotExist(err) {
			fmt.Printf("%v file does not exist\n", sourceFile)
			os.Exit(1)
		}

		targetFolder := filepath.Join(filepath.Dir(sourceFile), `jxl\`)
		if !exists(targetFolder) {
			err := os.Mkdir(targetFolder, 0755)
			if err != nil {
				fmt.Errorf("Error: %v", err)
			}
		}

		fmt.Println("Convert", sourceFile)
		fmt.Println("Size", ByteCountSI(fi.Size()))
		fmt.Println("Target Folder", targetFolder)

		input := sourceFile
		if !compatible {
			fmt.Println("File is not native supports by cjxl.exe, it will be convert to png")
			input = convertToPng(sourceFile, rootDir)
			input, _ = filepath.Abs(input)
		}

		output := filepath.Join(targetFolder, fmt.Sprintf("%v.jxl", filepath.Base(sourceFile)))
		if exists(output) {
			output = filepath.Join(targetFolder, fmt.Sprintf("%v_%v.jxl", filepath.Base(sourceFile), time.Now().Unix()))
		}

		convertToJxl(input, output, rootDir)

		if !compatible {
			fmt.Println("Remove temp file", input)
			os.Remove(input)
		}
	}
}

func convertToJxl(input, output, rootDir string) {

	fmt.Println("Convert to jxl in ", input)
	fmt.Println("Convert to jxl out", output)

	cjxlTempPath := filepath.Join(rootDir, "temp", "cjxl.exe")
	fi, err := os.Stat(cjxlTempPath)
	if os.IsNotExist(err) || fi.Size() != int64(len(magick_executable_data)) {
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

func exists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func convertToPng(source, rootDir string) string {
	magickTempPath := filepath.Join(rootDir, "temp", "magick.exe")
	fi, err := os.Stat(magickTempPath)
	if os.IsNotExist(err) || fi.Size() != int64(len(magick_executable_data)) {
		fmt.Println("Writing magick_temp.exe")
		os.WriteFile(magickTempPath, magick_executable_data, 0644)
	}

	ext := ".png"
	newFile := "temp_" + uuid.New().String() + ext
	newFile = filepath.Join(rootDir, "temp", newFile)
	start := time.Now()
	cmd := exec.Command(magickTempPath, "convert", source, newFile)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Convert to png in", time.Since(start), "to", newFile)

	files, err := filepath.Glob(strings.TrimSuffix(newFile, ext) + "*")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	for _, f := range files[1:] {
		os.Remove(f)
	}

	return files[0]
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
