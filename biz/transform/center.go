package transform

import (
	"encoding/json"
	"errors"
	"unsafe"
)

type Center struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Zoom      int64   `json:"zoom"`
}

var UniverseCenter = &Center{
	Latitude:  404,
	Longitude: 404,
	Zoom:      0,
}

func StringToCenter(str string) (*Center, error) {
	if str != "" {
		var center Center
		err := json.Unmarshal([]byte(str), &center)
		if err != nil {
			return nil, errors.New("StringToCenter failed")
		}
		return &center, nil
	}
	return UniverseCenter, nil
}

func PackCenter(center *Center) *string {
	bytesData, _ := json.Marshal(*center)
	return (*string)(unsafe.Pointer(&bytesData))
}
