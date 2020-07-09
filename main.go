package main

import (
	"flvParse/flv"
	"fmt"
	"os"
)

func main() {

	flvFile, err := os.Open("./test.flv")
	if err != nil {
		fmt.Printf("os.Open failed, err:%v\n", err)
	}

	buf := make([]byte, 0)
	f := new(flv.Flv)
	for true {

		tmpBuf := make([]byte, 1024)
		length, err := flvFile.Read(tmpBuf)
		if err != nil {
			fmt.Printf("flvFile.Read failed, err:%v\n", err)
		}
		fmt.Printf("read length:%v\n", length)
		if length == 0 {
			break
		}
		buf = append(buf, tmpBuf...)

		buf, err = f.Parse(buf)
		if err != nil {
			fmt.Printf("f.Parse failed, err:%v\n", err)
			os.Exit(0)
		}
		break
	}
}
