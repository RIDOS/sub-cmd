package cmd

import (
	"fmt"
	"io"
	"os"
)

type Output interface {
	Write(data []byte) error
}

type ConsoleOutput struct {
	wirter io.Writer
}

func (c *ConsoleOutput) Write(data []byte) error {
	_, err := fmt.Fprintf(c.wirter, "%s\n", data)
	if err != nil {
		return err
	}
	return nil
}

type FileOutput struct {
	filePath string
}

func (fileOutput *FileOutput) Write(data []byte) error {
	err := os.WriteFile(fileOutput.filePath+outputFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
