package util

import (
	"encoding/json"
	"fmt"
)

// LogJSON print object as json string on console
func LogJSON(value interface{}) {
	jb, _ := json.Marshal(value)
	fmt.Println(string(jb))
}
