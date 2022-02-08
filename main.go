package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dhcgn/jxldxoconverter/jxlhandler"
	"github.com/dhcgn/jxldxoconverter/magickhandler"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("context", "main")

	//go:embed usage.txt
	usage string
)

func main() {
	rootDir := filepath.Dir(os.Args[0])
	setupLogger(rootDir)

	if len(os.Args) == 1 || len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println(usage)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		log.Error("Need at least one Argument with a full path to a image file")
		fmt.Println(usage)
		os.Exit(1)
	}

	workingDir := createWorkingDir(rootDir)
	sourceFiles := os.Args[1:]

	for _, sourceFile := range sourceFiles {
		convertFile(sourceFile, workingDir)
	}
}

func convertFile(sourceFile, workingDir string) {
	fi, err := os.Stat(sourceFile)
	if os.IsNotExist(err) {
		log.Printf("%v file does not exist\n", sourceFile)
		return
	}

	targetFolder := createTargetFolder(sourceFile)

	log.WithFields(logrus.Fields{"source": sourceFile, "target": targetFolder, "size": ByteCountSI(fi.Size())}).Info("Converting")

	input := sourceFile
	compatible := jxlhandler.IsCompatible(sourceFile)
	if !compatible {
		log.Println("File is not native supports by cjxl.exe, it will be convert to png")
		input = magickhandler.ConvertToPng(sourceFile, workingDir, log.WithField("context", "magick"))
		input, _ = filepath.Abs(input)
	}

	output := filepath.Join(targetFolder, fmt.Sprintf("%v.jxl", filepath.Base(sourceFile)))
	if exists(output) {
		output = filepath.Join(targetFolder, fmt.Sprintf("%v_%v.jxl", filepath.Base(sourceFile), time.Now().Unix()))
	}

	jxlhandler.ConvertToJxl(input, output, workingDir, log.WithField("context", "jxl"))

	fiNew, err := os.Stat(output)
	if os.IsNotExist(err) {
		log.Errorf("%v file does not exist\n", fiNew)
		return
	}

	log.WithFields(logrus.Fields{"target": output, "new_size": ByteCountSI(fiNew.Size()), "saved": ByteCountSI(fi.Size() - fiNew.Size())}).Info("Converted")

	if !compatible {
		l := log.WithFields(logrus.Fields{"temp_file": input})
		l.Println("Removing")
		err := os.Remove(input)
		if err != nil {
			l.Errorf("Error: %v", err)
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
			log.Println("Error: ", err)
		}
	}
	return workingDir
}

func createTargetFolder(sourceFile string) string {
	targetFolder := filepath.Join(filepath.Dir(sourceFile), `jxl\`)
	if !exists(targetFolder) {
		err := os.Mkdir(targetFolder, 0755)
		if err != nil {
			log.Errorf("Error: %v", err)
		}
	}
	return targetFolder
}

func setupLogger(rootDir string) {
	filename := fmt.Sprintf("%v.log", time.Now().Format("2006-01-02"))
	filename = filepath.Join(rootDir, filename)
	fmt.Println("Logging to", filename)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		mw := io.MultiWriter(colorable.NewColorableStdout(), file)
		logrus.SetOutput(mw)
	} else {
		logrus.SetOutput(colorable.NewColorableStdout())
		log.Error(err)
	}
}
