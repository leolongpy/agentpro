package settings

import "fmt"

const VERSION = "0.2.8"

func GetVersion() string {
	fmt.Printf("build ver:\t%s\n", VERSION)
	return VERSION
}
