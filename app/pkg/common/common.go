package common

import "os"

const (
	Data = "./tne-helper-data/"
)

func init() {
	if err := os.MkdirAll(Data, 0755); err != nil {
		panic(err.Error())
	}
}
