package json

import "github.com/bytedance/sonic"

func Marshal(obj any) ([]byte, error) {
	return sonic.Marshal(obj)
}

func Unmarshal(data []byte, obj any) error {
	return sonic.Unmarshal(data, obj)
}
