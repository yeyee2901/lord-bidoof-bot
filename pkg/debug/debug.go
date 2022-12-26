package debug

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func DebugStruct(obj any) {
	if b, err := json.MarshalIndent(obj, "", "  "); err == nil {
		fmt.Println(reflect.TypeOf(obj).String(), string(b))
	} else {
		fmt.Println("Cannot debug:", reflect.TypeOf(obj))
	}
}
