package flv

import (
	"flvParse/util"
	"fmt"
)

const (
	Header = 0
	PreviousTagSize
	Tag
)
const (
	TypeFlagsReserved0Mark byte = 0b11111000
	TypeFlagsAudioMark     byte = 0b00000100
	TypeFlagsReserved1Mark byte = 0b00000010
	TypeFlagsVideoMark     byte = 0b00000001

	TagReservedMark byte = 0b11000000
	TagFilterMark   byte = 0b00100000
	TagTagTypeMark  byte = 0b00011111
)

var TagTypeMap map[uint8]string = map[uint8]string{
	8:  "audio",
	9:  "video",
	18: "script data",
}

type Flv struct {
	State              int
	PreviousTagSizeNum int
}

func (f *Flv) Parse(buf []byte) ([]byte, error) {
	//fmt.Println("Parse")
	var err error
	var ok bool
	if f.State == Header {
		buf, ok, err = f.parseHeader(buf)
		if ok {
			f.State = PreviousTagSize
		}
	}
	if f.State == PreviousTagSize {
		buf, ok, err = f.parsePreviousTagSize(buf)
		if ok {
			f.State = Tag
			f.PreviousTagSizeNum++
		}
	}
	if f.State == Tag {
		buf, ok, err = f.parseTag(buf)
		if ok {
			f.State = PreviousTagSize
		}
	}
	if err != nil {
		return nil, fmt.Errorf("f.parse failed, err:%v", err)
	}
	return buf, nil
}

func (f *Flv) parseHeader(buf []byte) ([]byte, bool, error) {
	fmt.Println("parseHeader")
	if len(buf) < 9 {
		return buf, false, nil
	}
	if buf[0] != 0x46 {
		return nil, false, fmt.Errorf("signature0 != 0x46, signature0:%x", buf[0])
	}
	fmt.Println("Signature0 is 0x46")

	if buf[1] != 0x4C {
		return nil, false, fmt.Errorf("signature1 != 0x4C, signature1:%x", buf[1])
	}
	fmt.Println("Signature1 is 0x4C")

	if buf[2] != 0x56 {
		return nil, false, fmt.Errorf("signature2 != 0x56, signature2:%x", buf[2])
	}
	fmt.Println("Signature2 is 0x56")

	if buf[3] != 0x01 {
		return nil, false, fmt.Errorf("version != 0x01, version:%x", buf[3])
	}
	fmt.Println("Version is 0x01")

	typeFlagsReserved0 := (buf[4] & TypeFlagsReserved0Mark) >> 3
	if typeFlagsReserved0 != 0 {
		return nil, false, fmt.Errorf("TypeFlagsReserved0 != 0, TypeFlagsReserved:%x", typeFlagsReserved0)
	}
	fmt.Println("TypeFlagsReserved0 is 0")

	typeFlagsAudio := (buf[4] & TypeFlagsAudioMark) >> 2
	fmt.Printf("TypeFlagsAudio is %v\n", typeFlagsAudio)

	typeFlagsReserved1 := (buf[4] & TypeFlagsReserved1Mark) >> 1
	if typeFlagsReserved1 != 0 {
		return nil, false, fmt.Errorf("typeFlagsReserved1 != 0, TypeFlagsReserved1:%x", typeFlagsReserved1)
	}
	fmt.Printf("TypeFlagsReserved1 is %v\n", typeFlagsReserved1)

	typeFlagsVideo := (buf[4] & TypeFlagsVideoMark) >> 0
	fmt.Printf("TypeFlagsVideo is %v\n", typeFlagsVideo)

	DataOffset, err := util.BytesToUint32ByBigEndian(buf[5:9])
	if err != nil {
		return nil, false, fmt.Errorf("util.BytesToUint32ByBigEndian failed, err:%v", err)
	}
	if DataOffset != 9 {
		return nil, false, fmt.Errorf("DataOffset != 9, DataOffset != 9:%v", DataOffset)
	}
	fmt.Println("DataOffset is 9")

	return buf[9:], true, nil
}

func (f *Flv) parsePreviousTagSize(buf []byte) ([]byte, bool, error) {
	if len(buf) < 4 {
		return nil, false, nil
	}

	previousTagSize, err := util.BytesToUint32ByBigEndian(buf[:4])
	if err != nil {
		return nil, false, fmt.Errorf("util.BytesToUint32ByBigEndian(buf[:4]) failed, err:%v", err)
	}
	fmt.Printf("PreviousTagSize%v is %v\n", f.PreviousTagSizeNum, previousTagSize)

	return buf[4:], true, nil
}

func (f *Flv) parseTag(buf []byte) ([]byte, bool, error) {

	if len(buf) < 11 {
		return buf, false, nil
	}
	dataSize, err := util.BytesToUint32ByBigEndian(buf[1:4])
	//fmt.Printf("DataSize：%v buf:%x\n", dataSize, buf[1:4])
	if err != nil {
		return nil, false, fmt.Errorf("util.BytesToUint32ByBigEndian failed, err:%v", err)
	}
	if len(buf) < int(11+dataSize) {
		return buf, false, nil
	}

	index := 0
	reserved := buf[index] & TagReservedMark >> 6
	if reserved != 0 {
		return nil, false, fmt.Errorf("reserved != 0, reserved:%v", reserved)
	}
	fmt.Println("Reserved is 0")

	filter := buf[index] & TagFilterMark >> 5
	fmt.Printf("Filter is %v\n", filter)

	TagType, err := util.BytesToUint8ByBigEndian(buf[index] & TagTagTypeMark)
	if err != nil {
		return nil, false, fmt.Errorf("util.BytesToUint8ByBigEndian failed, err:%v", err)
	}
	if TagType != 8 && TagType != 9 && TagType != 18 {
		return nil, false, fmt.Errorf("TagType != 8 && TagType != 9 && TagType != 18, TagType:%v", TagType)
	}
	fmt.Printf("TagType is %v\n", TagTypeMap[TagType])
	index += 1

	fmt.Printf("DataSize is %v\n", dataSize)
	index += 3

	timestamp, err := util.BytesToUint32ByBigEndian(buf[index : index+3])
	if err != nil {
		return nil, false, fmt.Errorf("util.BytesToUint32ByBigEndian failed, err:%v", err)
	}
	fmt.Printf("Timestamp is %v\n", timestamp)
	index += 3

	timestampExtended, err := util.BytesToUint8ByBigEndian(buf[index])
	if err != nil {
		return nil, false, fmt.Errorf("util.BytesToUint8ByBigEndian failed, err:%v", err)
	}
	fmt.Printf("TimestampExtended is %v\n", timestampExtended)
	index += 1

	streamID, err := util.BytesToUint32ByBigEndian(buf[index : index+3])
	if err != nil {
		return nil, false, fmt.Errorf("util.BytesToUint32ByBigEndian failed, err:%v", err)
	}
	if streamID != 0 {
		return nil, false, fmt.Errorf("streamID != 0, streamID:%v", streamID)
	}
	fmt.Println("streamID is 0")
	index += 3

	if TagType == 8 {
		length, err := f.parseAudioTagHeader(buf[index:])
		if err != nil {
			return nil, false, fmt.Errorf("f.parseAudioTagHeader failed, err:%v", err)
		}
		index += length
	}


	return buf[:], true, nil
}

func (f *Flv) parseAudioTagHeader(buf []byte) (int, error) {
	return 0, nil
}

func (f *Flv) parseVideoTagHeader(buf []byte) (int, error) {
	return 0, nil
}

func (f *Flv) parseEncryptionHeader(buf []byte) (int, error) {
	return 0, nil
}

func (f *Flv) parseFilterParams(buf []byte) (int, error) {
	return 0, nil
}

func (f *Flv) parseData(buf []byte) (int, error) {
	return 0, nil
}
