package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dhcgn/jxldxoconverter/jxlhandler"
	"github.com/dhcgn/jxldxoconverter/magickhandler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Need at least one Argument with a full path to a .tif file")
		os.Exit(1)
	}

	rootDir := filepath.Dir(os.Args[0])
	workingDir := createWorkingDir(rootDir)
	sourceFiles := os.Args[1:]

	for _, sourceFile := range sourceFiles {

		fi, err := os.Stat(sourceFile)
		if os.IsNotExist(err) {
			fmt.Printf("%v file does not exist\n", sourceFile)
			continue
		}

		targetFolder := createTargetFolder(sourceFile)

		fmt.Println("Source", sourceFile)
		fmt.Println("Size", ByteCountSI(fi.Size()))
		fmt.Println("Target Folder", targetFolder)

		input := sourceFile
		compatible := jxlhandler.IsCompatible(sourceFile)
		if !compatible {
			fmt.Println("File is not native supports by cjxl.exe, it will be convert to png")
			input = magickhandler.ConvertToPng(sourceFile, workingDir)
			input, _ = filepath.Abs(input)
		}

		output := filepath.Join(targetFolder, fmt.Sprintf("%v.jxl", filepath.Base(sourceFile)))
		if exists(output) {
			output = filepath.Join(targetFolder, fmt.Sprintf("%v_%v.jxl", filepath.Base(sourceFile), time.Now().Unix()))
		}

		jxlhandler.ConvertToJxl(input, output, workingDir)

		fiNew, err := os.Stat(output)
		if os.IsNotExist(err) {
			fmt.Printf("%v file does not exist\n", fiNew)
			continue
		}

		fmt.Println("New Size", ByteCountSI(fiNew.Size()), "Diff", ByteCountSI(fi.Size()-fiNew.Size()))

		if !compatible {
			fmt.Println("Remove temp file", input)
			os.Remove(input)
		}
	}
}

func exists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
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

func createWorkingDir(rootDir string) string {
	workingDir := filepath.Join(rootDir, "temp")
	if !exists(workingDir) {
		err := os.Mkdir(workingDir, 0755)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
	return workingDir
}

func createTargetFolder(sourceFile string) string {
	targetFolder := filepath.Join(filepath.Dir(sourceFile), `jxl\`)
	if !exists(targetFolder) {
		err := os.Mkdir(targetFolder, 0755)
		if err != nil {
			fmt.Errorf("Error: %v", err)
		}
	}
	return targetFolder
}
