package helper

import (
	"os"
)

func LoadFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func SaveFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
