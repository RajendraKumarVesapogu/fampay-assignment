package utils

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"
)

func GetPaginationOffset(page int, size int) int {
	return (page - 1) * size
}

func GetContextWithTimeout(ctx context.Context, timeout time.Duration) context.Context {
	Ctx, cancel := context.WithCancel(ctx)

	timer := time.NewTimer(timeout)

	go func() {
		defer cancel()
		<-timer.C
	}()

	return Ctx
}

func EncodeToGob(data any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DecodeFromGob(data []byte, result any) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(result)
	if err != nil {
		return err
	}
	return nil
}