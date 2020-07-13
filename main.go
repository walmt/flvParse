package main

import (
	"flvParse/flv"
	"fmt"
	"io"
	"os"
)

func main() {

	flvFile, err := os.Open("./test.flv")
	if err != nil {
		fmt.Printf("os.Open(\"./test.flv\") failed, err:%v\n", err)
	}

	buf := make([]byte, 0)
	f := new(flv.Flv)

	//nums := 0
	//times := 1

	for true {

		tmpBuf := make([]byte, 1024000)
		length, errRead := flvFile.Read(tmpBuf)
		if errRead != nil && errRead != io.EOF {
			fmt.Printf("flvFile.Read failed, err:%v\n", errRead)
		}
		if length == 0 {
			fmt.Printf("read end\n")
			break
		}
		fmt.Printf("read length:%v\n", length)
		buf = append(buf, tmpBuf[:length]...)

		buf, err = f.Parse(buf)
		if err != nil {
			fmt.Printf("f.Parse failed, err:%v\n", err)
			os.Exit(-1)
		}
		if errRead == io.EOF {
			fmt.Printf("already read and deal")
			os.Exit(0)
		}

		//nums++
		//if nums == times {
		//	break
		//}
	}

	fmt.Println()
}
