package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func BytesToUint32ByBigEndian(buf []byte) (uint32, error) {

	if len(buf) > 4 {
		return 0, fmt.Errorf("buf len is not 4, len:%v", len(buf))
	}

	b := make([]byte, 4)
	less := 4 - len(buf)
	for i := 0; i < len(buf); i++ {
		b[i+less] = buf[i]
	}

	bytesBuffer := bytes.NewBuffer(b)

	var x uint32
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); err != nil {
		return 0, fmt.Errorf("binary.Read failed, err:%v", err)
	}

	return x, nil
}

func BytesToUint8ByBigEndian(buf byte) (uint8, error) {

	b := make([]byte, 1)
	b[0] = buf

	bytesBuffer := bytes.NewBuffer(b)

	var x uint8
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); err != nil {
		return 0, fmt.Errorf("binary.Read failed, err:%v", err)
	}

	return x, nil
}
