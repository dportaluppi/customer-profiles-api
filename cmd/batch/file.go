package main

import (
	"encoding/json"
	"os"
)

type File struct {
	file *os.File
}

func CreateFile(name string) (*File, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	return &File{
		file: file,
	}, nil
}

func (f *File) Start() error {
	_, err := f.file.WriteString("[\n")

	return err
}

func (f *File) Append(data any) error {
	jsonData, err := f.marshalData(data)
	if err != nil {
		return err
	}

	_, err = f.file.WriteString(jsonData + ",\n")

	return err
}

func (f *File) Finish(data any) error {
	jsonData, err := f.marshalData(data)
	if err != nil {
		return err
	}

	_, err = f.file.WriteString(jsonData + "]")
	if err != nil {
		return err
	}

	return f.file.Close()
}

func (f *File) marshalData(data any) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
