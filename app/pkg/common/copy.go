package common

import "os"

func CopyFile(dst, src string) {
	if data, err := os.ReadFile(src); err != nil {
		panic(err.Error())
	} else if err := os.WriteFile(dst, data, 0644); err != nil {
		panic(err.Error())
	}
}
