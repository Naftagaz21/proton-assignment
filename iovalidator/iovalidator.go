package iovalidator

import (
	"errors"
	"os"
	"path/filepath"
)

type Validator interface {
	ValidateAndProcessInput(input string) ([]string, error)
	ValidateOutput(output string) error
}

type validator struct{}

func (*validator) ValidateAndProcessInput(inputPath string) ([]string, error) {
	fileInfo, err := os.Stat(inputPath)

	if err != nil && os.IsNotExist(err) {
		return nil, errors.New("Input directory does not exists")
	}

	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, errors.New("Specified input argument is not a directory")
	}

	// Does not search through the entire directory tree
	matches, err := filepath.Glob(inputPath + "/*.md")
	if err != nil {
		return nil, err
	}

	if len(matches) < 1 {
		return nil, errors.New("No markdown files found in input folder")
	}

	return matches, nil
}

func (*validator) ValidateOutput(outputPath string) error {
	fileInfo, err := os.Stat(outputPath)
	if err == nil {
		if !fileInfo.IsDir() {
			return errors.New("The specified output path is a file. Please specify a folder location instead")
		}
		return nil
	}

	if os.IsNotExist(err) {
		err := os.MkdirAll(outputPath, 0755)
		if err != nil {
			return errors.New("Unable to create output folder")
		}
		return nil
	}

	return err
}

func NewValidator() Validator {
	return &validator{}
}
