package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dhcgn/jxldxoconverter/config"
	"github.com/dhcgn/jxldxoconverter/helper"
	"github.com/dhcgn/jxldxoconverter/jxlhandler"
	"github.com/dhcgn/jxldxoconverter/magickhandler"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("context", "main")

	//go:embed usage.txt
	usage string

	version = "0.0.1"
)

func main() {
	fmt.Println("jxl dxo converter", version)

	rootDir := filepath.Dir(os.Args[0])
	setupLogger(rootDir)
	c := config.GetConfig(rootDir)

	if len(os.Args) == 1 || len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println(usage)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		log.Error("Need at least one Argument with a full path to a image file")
		fmt.Println(usage)
		os.Exit(1)
	}

	log.WithFields(logrus.Fields{"images": len(os.Args) - 1}).Info("Started")

	workingDir := createWorkingDir(rootDir)
	sourceFiles := os.Args[1:]

	for _, sourceFile := range sourceFiles {
		ffs := c.GetFileFormatSetting(sourceFile)
		convertFile(sourceFile, workingDir, ffs)
	}
}

func convertFile(sourceFile, workingDir string, ffs config.FileFormatSetting) {
	fi, err := os.Stat(sourceFile)
	if os.IsNotExist(err) {
		log.Printf("%v file does not exist\n", sourceFile)
		return
	}

	targetFolder := createTargetFolder(sourceFile)

	log.WithFields(logrus.Fields{"source": sourceFile, "target": targetFolder, "size": helper.ByteCountSI(fi.Size())}).Info("Converting")

	input := sourceFile
	compatible := jxlhandler.IsCompatible(sourceFile)
	if !compatible {
		log.Println("File is not native supports by cjxl.exe, it will be convert to png")
		input = magickhandler.ConvertToPng(sourceFile, workingDir, log.WithField("context", "magick"))
		input, _ = filepath.Abs(input)
	}

	output := filepath.Join(targetFolder, fmt.Sprintf("%v.jxl", filepath.Base(sourceFile)))
	if helper.Exists(output) {
		output = filepath.Join(targetFolder, fmt.Sprintf("%v_%v.jxl", filepath.Base(sourceFile), time.Now().Unix()))
	}

	jxlhandler.ConvertToJxl(input, output, workingDir, ffs, log.WithField("context", "jxl"))

	fiNew, err := os.Stat(output)
	if os.IsNotExist(err) {
		log.Errorf("%v file does not exist\n", fiNew)
		return
	}

	log.WithFields(logrus.Fields{"target": output, "new_size": helper.ByteCountSI(fiNew.Size()), "saved": helper.ByteCountSI(fi.Size() - fiNew.Size())}).Info("Converted")

	if !compatible {
		l := log.WithFields(logrus.Fields{"temp_file": input})
		l.Println("Removing")
		err := os.Remove(input)
		if err != nil {
			l.Errorf("Error: %v", err)
		}
	}

	if !ffs.DefaultConfig && ffs.DeleteSourceFile {
		l := log.WithFields(logrus.Fields{"input_file": sourceFile})
		l.Println("Removing")
		err := os.Remove(sourceFile)
		if err != nil {
			l.Errorf("Error: %v", err)
		}
	}
}

func createWorkingDir(rootDir string) string {
	workingDir := filepath.Join(rootDir, "temp")
	if !helper.Exists(workingDir) {
		err := os.Mkdir(workingDir, 0755)
		if err != nil {
			log.Println("Error: ", err)
		}
	}
	return workingDir
}

func createTargetFolder(sourceFile string) string {
	targetFolder := filepath.Join(filepath.Dir(sourceFile), `jxl\`)
	if !helper.Exists(targetFolder) {
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
