package util

import "os"
import "strings"

func ExpandPath(path string) (string, error) {
	path = os.ExpandEnv(path)
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = home + path[1:]
	}
	return path, nil
}
