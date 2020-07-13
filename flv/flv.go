package flv

import (
	"flvParse/util"
	"fmt"
	"os"
)

const (
	Header          = 0
	PreviousTagSize = 1
	Tag             = 2
)

const (
	TypeFlagsReserved0Mark byte = 0b11111000
	TypeFlagsAudioMark     byte = 0b00000100
	TypeFlagsReserved1Mark byte = 0b00000010
	TypeFlagsVideoMark     byte = 0b00000001

	TagReservedMark byte = 0b11000000
	TagFilterMark   byte = 0b00100000
	TagTagTypeMark  byte = 0b00011111

	FrameTypeMark byte = 0b11110000
	CodecIDMark   byte = 0b00001111

	SoundFormatMark byte = 0b11110000
	SoundRateMark   byte = 0b00001100
	SoundSizeMark   byte = 0b00000010
	SoundTypeMark   byte = 0b00000001
)

const (
	TagTypeAudio      = 8
	TagTypeVideo      = 9
	TagTypeScriptData = 18

	FilterNoPreProcessingRequired = 0
	FilterPreProcessing           = 1

	FrameTypeKeyFrame                = 1 // for AVC, a seekable frame
	FrameTypeInterFrame              = 2 // for AVC, a non-seekable frame
	FrameTypeDisposableInterFrame    = 3 // H.263 only
	FrameTypeGeneratedKeyFrame       = 4 // reserved for server use only
	FrameTypeVideoInfoOrCommandFrame = 5

	CodecIDSorensonH263           = 2
	CodecIDScreenVideo            = 3
	CodecIDOn2Vp6                 = 4
	CodecIDOn2Vp6WithAlphaChannel = 5
	CodecIDScreenVideoVersion2    = 6
	CodecIDAvc                    = 7

	AvcPacketTypeAvcSequenceHeader = 0
	AvcPacketTypeAvcNalu           = 1
	AvcPacketTypeAvcEndOfSequence  = 2 // lower level NALU sequence ender is not required or supported

	ScriptDataValueTypeNumber          = 0
	ScriptDataValueTypeBoolean         = 1
	ScriptDataValueTypeString          = 2
	ScriptDataValueTypeObject          = 3
	ScriptDataValueTypeMovieClip       = 4
	ScriptDataValueTypeNull            = 5
	ScriptDataValueTypeUndefined       = 6
	ScriptDataValueTypeReference       = 7
	ScriptDataValueTypeEcmaArray       = 8
	ScriptDataValueTypeObjectEndMarker = 9
	ScriptDataValueTypeStrictArray     = 10
	ScriptDataValueTypeDate            = 11
	ScriptDataValueTypeLongString      = 12

	SoundFormatLinearPcmPlatformEndian = 0
	SoundFormatAdpcm                   = 1
	SoundFormatMp3                     = 2
	SoundFormatLinearPcmLittleEndian   = 3
	SoundFormatNellymoser16kHzMono     = 4
	SoundFormatNellymoser8kHzMono      = 5
	SoundFormatNellymoser              = 6
	SoundFormatG711ALawLogarithmicPcm  = 7
	SoundFormatG711MuLawLogarithmicPcm = 8
	SoundFormatreserved                = 9
	SoundFormatAAC                     = 10
	SoundFormatSpeex                   = 11
	SoundFormatMP3_8kHz                = 14
	SoundFormatDeviceSpecificSound     = 15

	SoundRate5_5kHz = 0
	SoundRate11kHz  = 1
	SoundRate22kHz  = 2
	SoundRate44kHz  = 3

	SoundSize8BitSamples  = 0
	SoundSize16BitSamples = 1

	soundTypeMonoSound   = 0
	soundTypeStereoSound = 1

	AACPacketTypeAacSequenceHeader = 0
	AACPacketTypeAacRaw            = 1
)

var TagTypeMap = map[uint8]string{
	TagTypeAudio:      "audio",
	TagTypeVideo:      "video",
	TagTypeScriptData: "script data",
}

var FilterMap = map[uint8]string{
	FilterNoPreProcessingRequired: "No pre-processing required",
	FilterPreProcessing:           "Pre-processing",
}

var FrameTypeMap = map[uint8]string{
	FrameTypeKeyFrame:                "key frame",
	FrameTypeInterFrame:              "inter frame",
	FrameTypeDisposableInterFrame:    "disposable inter frame",
	FrameTypeGeneratedKeyFrame:       "generated key frame",
	FrameTypeVideoInfoOrCommandFrame: "video info/command frame",
}

var CodeIdMap = map[uint8]string{
	CodecIDSorensonH263:           "Sorenson H.263",
	CodecIDScreenVideo:            "Screen video",
	CodecIDOn2Vp6:                 "On2 VP6",
	CodecIDOn2Vp6WithAlphaChannel: "On2 VP6 with alpha channel",
	CodecIDScreenVideoVersion2:    "Screen video version 2",
	CodecIDAvc:                    "AVC",
}

var AvcPacketTypeMap = map[uint8]string{
	AvcPacketTypeAvcSequenceHeader: "AVC sequence header",
	AvcPacketTypeAvcNalu:           "AVC NALU",
	AvcPacketTypeAvcEndOfSequence:  "AVC end of sequence",
}

var SoundFormatMap = map[uint8]string{
	SoundFormatLinearPcmPlatformEndian: "Linear PCM, platform endian",
	SoundFormatAdpcm:                   "ADPCM",
	SoundFormatMp3:                     "MP3",
	SoundFormatLinearPcmLittleEndian:   "Linear PCM, little endian",
	SoundFormatNellymoser16kHzMono:     "Nellymoser 16 kHz mono",
	SoundFormatNellymoser8kHzMono:      "Nellymoser 8 kHz mono",
	SoundFormatNellymoser:              "Nellymoser",
	SoundFormatG711ALawLogarithmicPcm:  "G.711 A-law logarithmic PCM",
	SoundFormatG711MuLawLogarithmicPcm: "G.711 mu-law logarithmic PCM",
	SoundFormatreserved:                "reserved",
	SoundFormatAAC:                     "AAC",
	SoundFormatSpeex:                   "Speex",
	SoundFormatMP3_8kHz:                "MP3 8 kHz",
	SoundFormatDeviceSpecificSound:     "Device-specific sound",
}

var SoundRateMap = map[uint8]string{
	SoundRate5_5kHz: "5.5 kHz",
	SoundRate11kHz:  "11 kHz",
	SoundRate22kHz:  "22 kHz",
	SoundRate44kHz:  "44 kHz",
}

var SoundSizeMap = map[uint8]string{
	SoundSize8BitSamples:  "8-bit samples",
	SoundSize16BitSamples: "16-bit samples",
}

var SoundTypeMap = map[uint8]string{
	soundTypeMonoSound:   "Mono sound",
	soundTypeStereoSound: "Stereo sound",
}

var AACPacketTypeMap = map[uint8]string{
	AACPacketTypeAacSequenceHeader: "AAC sequence header",
	AACPacketTypeAacRaw:            "AAC raw",
}

var ScriptDataValueTypeSet = map[uint8]string{
	ScriptDataValueTypeNumber:          "Number",
	ScriptDataValueTypeBoolean:         "Boolean",
	ScriptDataValueTypeString:          "String",
	ScriptDataValueTypeObject:          "Object",
	ScriptDataValueTypeMovieClip:       "MovieClip (reserved, not supported)",
	ScriptDataValueTypeNull:            "Null",
	ScriptDataValueTypeUndefined:       "Undefined",
	ScriptDataValueTypeReference:       "Reference",
	ScriptDataValueTypeEcmaArray:       "ECMA array",
	ScriptDataValueTypeObjectEndMarker: "Object end marker",
	ScriptDataValueTypeStrictArray:     "Strict array",
	ScriptDataValueTypeDate:            "Date",
	ScriptDataValueTypeLongString:      "Long string",
}

type Flv struct {
	State              int
	PreviousTagSizeNum int
	CurrentTag         *CurrentTag
}

type CurrentTag struct {
	Length        int
	Filter        uint8
	TagType       uint8
	FrameType     uint8
	CodeId        uint8
	AVCPacketType uint8
	SoundFormat   uint8
	AACPacketType uint8
	SoundRate     uint8
}

var h264File *os.File
var aacFile *os.File

func (f *Flv) Parse(buf []byte) ([]byte, error) {
	//fmt.Println("Parse")

	var ok bool
	var err error

	if h264File == nil {
		h264File, err = os.OpenFile("./test.h264", os.O_CREATE|os.O_RDWR, 0)
		if err != nil {
			return nil, fmt.Errorf("os.OpenFile(\"./test.h264\", os.O_CREATE|os.O_RDWR, 0) failed, err:%v\n", err)
		}
	}

	if aacFile == nil {
		aacFile, err = os.OpenFile("./test.aac", os.O_CREATE|os.O_RDWR, 0)
		if err != nil {
			return nil, fmt.Errorf("os.OpenFile(\"./test.aac\", os.O_CREATE|os.O_RDWR, 0) failed, err:%v\n", err)
		}
	}
	for true {
		if f.State == Header {
			buf, ok, err = f.parseHeader(buf)
			if err != nil {
				return nil, fmt.Errorf("f.parseHeader failed, err:%v", err)
			}
			if ok {
				f.State = PreviousTagSize
				fmt.Println()
			}
		}
		if f.State == PreviousTagSize {
			buf, ok, err = f.parsePreviousTagSize(buf)
			if err != nil {
				return nil, fmt.Errorf("f.parsePreviousTagSize failed, err:%v", err)
			}
			if ok {
				f.State = Tag
				f.PreviousTagSizeNum++
				fmt.Println()
			}
		}
		if f.State == Tag {
			buf, ok, err = f.parseTag(buf)
			if err != nil {
				return nil, fmt.Errorf("f.parseTag failed, err:%v", err)
			}
			if ok {
				f.State = PreviousTagSize
				fmt.Println()
			}
		}
		if !ok || len(buf) == 0 {
			return buf, nil
		}
	}

	return nil, nil
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

	f.CurrentTag = new(CurrentTag)
	var index int

	if len(buf) < 11 {
		return buf, false, nil
	}
	dataSize, err := util.BytesToUint32ByBigEndian(buf[1:4])
	if err != nil {
		return nil, false, fmt.Errorf("util.BytesToUint32ByBigEndian failed, err:%v", err)
	}
	if len(buf) < int(11+dataSize) {
		return buf, false, nil
	}

	f.CurrentTag.Length = int(11 + dataSize)

	reserved := buf[index] & TagReservedMark >> 6
	if reserved != 0 {
		return nil, false, fmt.Errorf("reserved != 0, reserved:%v", reserved)
	}
	fmt.Println("Reserved is 0")

	filter := buf[index] & TagFilterMark >> 5
	filterString, ok := FilterMap[filter]
	if !ok {
		return nil, false, fmt.Errorf("FilterMap[filter] failed, filter:%v", filter)
	}
	f.CurrentTag.Filter = filter
	fmt.Printf("Filter is %v\n", filterString)

	tagType := util.BytesToUint8ByBigEndian(buf[index] & TagTagTypeMark)
	if _, ok := TagTypeMap[tagType]; !ok {
		return nil, false, fmt.Errorf("TagType is illegal, TagType:%v", tagType)
	}
	f.CurrentTag.TagType = tagType
	fmt.Printf("TagType is %v\n", TagTypeMap[tagType])
	index += 1

	fmt.Printf("DataSize is %v\n", dataSize)
	index += 3

	timestamp, err := util.BytesToUint32ByBigEndian(buf[index : index+3])
	if err != nil {
		return nil, false, fmt.Errorf("util.BytesToUint32ByBigEndian failed, err:%v", err)
	}
	fmt.Printf("Timestamp is %v\n", timestamp)
	index += 3

	timestampExtended := util.BytesToUint8ByBigEndian(buf[index])
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

	if f.CurrentTag.TagType == TagTypeAudio {
		index, err = f.parseAudioTagHeader(buf, index)
		if err != nil {
			return nil, false, fmt.Errorf("f.parseAudioTagHeader failed, err:%v", err)
		}
	}

	if f.CurrentTag.TagType == TagTypeVideo {
		index, err = f.parseVideoTagHeader(buf, index)
		if err != nil {
			return nil, false, fmt.Errorf("f.parseVideoTagHeader failed, err:%v", err)
		}
	}

	if f.CurrentTag.Filter == FilterPreProcessing {
		index, err = f.parseEncryptionHeader(buf, index)
		if err != nil {
			return nil, false, fmt.Errorf("f.parseEncryptionHeader failed, err:%v", err)
		}

		index, err = f.parseFilterParams(buf, index)
		if err != nil {
			return nil, false, fmt.Errorf("f.parseFilterParams failed, err:%v", err)
		}
	}

	index, err = f.parseData(buf, index)
	if err != nil {
		return nil, false, fmt.Errorf("f.parseData failed, err:%v", err)
	}

	//fmt.Printf("end index:%v\n", index)
	return buf[11+dataSize:], true, nil
}

func (f *Flv) parseAudioTagHeader(buf []byte, index int) (int, error) {
	if len(buf) < 1 {
		return 0, fmt.Errorf("len(buf) < 1")
	}

	soundFormat := util.BytesToUint8ByBigEndian((buf[index] & SoundFormatMark) >> 4)
	soundFormatString, ok := SoundFormatMap[soundFormat]
	if !ok {
		return 0, fmt.Errorf("SoundFormatMap[soundFormat] failed, soundFormat:%v", soundFormat)
	}
	f.CurrentTag.SoundFormat = soundFormat
	fmt.Printf("soundFormat is %v\n", soundFormatString)

	soundRate := util.BytesToUint8ByBigEndian((buf[index] & SoundRateMark) >> 2)
	soundRateString, ok := SoundRateMap[soundRate]
	if !ok {
		return 0, fmt.Errorf("SoundRateMap[soundRate] failed, soundRate:%v", soundRate)
	}
	fmt.Printf("soundRate is %v\n", soundRateString)

	soundSize := util.BytesToUint8ByBigEndian((buf[index] & SoundSizeMark) >> 1)
	soundSizeString, ok := SoundSizeMap[soundSize]
	if !ok {
		return 0, fmt.Errorf("SoundSizeMap[soundRate] failed, soundSize:%v", soundSize)
	}
	fmt.Printf("soundSize is %v\n", soundSizeString)

	soundType := util.BytesToUint8ByBigEndian((buf[index] & SoundTypeMark) >> 0)
	soundTypeString, ok := SoundTypeMap[soundType]
	if !ok {
		return 0, fmt.Errorf("soundTypeMap[soundType] failed, soundType:%v", soundType)
	}
	fmt.Printf("soundType is %v\n", soundTypeString)

	index += 1

	if f.CurrentTag.SoundFormat == SoundFormatAAC {
		if len(buf[index:]) < 1 {
			return 0, fmt.Errorf("len(buf[%v:]) < 1", index)
		}
		aacPacketType := buf[index]
		aacPacketTypeString, ok := AACPacketTypeMap[aacPacketType]
		if !ok {
			return 0, fmt.Errorf("AACPacketTypeMap[aacPacketType] failed, aacPacketType:%v", aacPacketType)
		}
		f.CurrentTag.AACPacketType = aacPacketType
		fmt.Printf("aacPacketType is %v\n", aacPacketTypeString)

		index += 1
	}

	return index, nil
}

func (f *Flv) parseVideoTagHeader(buf []byte, index int) (int, error) {
	if len(buf) < 1 {
		return 0, fmt.Errorf("len(buf) < 1")
	}

	frameType := util.BytesToUint8ByBigEndian(buf[index] & FrameTypeMark >> 4)
	frameTypeString, ok := FrameTypeMap[frameType]
	if !ok {
		return 0, fmt.Errorf("FrameTypeMap[frameType] is not ok, frameType:%v", frameType)
	}
	f.CurrentTag.FrameType = frameType
	fmt.Printf("FrameType is %v\n", frameTypeString)

	codeId := util.BytesToUint8ByBigEndian(buf[index] & CodecIDMark)
	codeIdString, ok := CodeIdMap[codeId]
	if !ok {
		return 0, fmt.Errorf("CodeIdMap[codeId] is not ok, codeId:%v", codeId)
	}
	f.CurrentTag.CodeId = codeId
	fmt.Printf("CodeId is %v\n", codeIdString)

	index += 1

	if codeId == CodecIDAvc {
		if len(buf[index:]) < 4 {
			return 0, fmt.Errorf("len(buf[index:]) < 4")
		}

		avcPacketType := util.BytesToUint8ByBigEndian(buf[index])
		avcPacketTypeString, ok := AvcPacketTypeMap[avcPacketType]
		if !ok {
			return 0, fmt.Errorf("AvcPacketTypeMap[avcPacketType] is not ok, avcPacketType:%v", avcPacketType)
		}
		f.CurrentTag.AVCPacketType = avcPacketType
		fmt.Printf("AvcPacketType is %v\n", avcPacketTypeString)
		index += 1

		compositionTime, err := util.BytesToInt32ByBigEndian(buf[index : index+3])
		if err != nil {
			return 0, fmt.Errorf("util.BytesToInt32ByBigEndian failed, err:%v", err)
		}
		if avcPacketType != AvcPacketTypeAvcNalu && compositionTime != 0 {
			return 0, fmt.Errorf("CompositionTime must to be 0")
		}
		fmt.Printf("CompositionTime is %v\n", compositionTime)
		index += 3
	}

	return index, nil
}

func (f *Flv) parseEncryptionHeader(buf []byte, index int) (int, error) {
	return 0, fmt.Errorf("parseEncryptionHeader error")
}

func (f *Flv) parseFilterParams(buf []byte, index int) (int, error) {
	return 0, fmt.Errorf("parseFilterParams error")
}

func (f *Flv) parseData(buf []byte, index int) (int, error) {

	var err error
	if f.CurrentTag.TagType == TagTypeAudio {
		index, err = f.parseAudioData(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseAudioData failed, err:%v", err)
		}
	}
	if f.CurrentTag.TagType == TagTypeVideo {
		index, err = f.parseVideoData(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseVideoData failed, err:%v", err)
		}
	}
	if f.CurrentTag.TagType == TagTypeScriptData {
		index, err = f.parseScriptData(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseScriptData failed, err:%v", err)
		}
	}
	return index, nil
}

func (f *Flv) parseAudioData(buf []byte, index int) (int, error) {

	var err error

	if f.CurrentTag.Filter == FilterPreProcessing {
		index, err = f.parseAudioDataEncryptedBody(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseAudioDataEncryptedBody failed, err:%v", err)
		}
	} else {
		index, err = f.parseAudioDataAudioTagBody(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseAudioDataAudioTagBody failed, err:%v", err)
		}
	}
	return index, nil
}

func (f *Flv) parseAudioDataEncryptedBody(buf []byte, index int) (int, error) {
	return 0, fmt.Errorf("parseAudioDataEncryptedBody error")
}

func (f *Flv) parseAudioDataAudioTagBody(buf []byte, index int) (int, error) {
	var err error
	if f.CurrentTag.SoundFormat == SoundFormatAAC {
		index, err = f.parseAacAudioData(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseAacAudioData failed, err:%v", err)
		}
	} else {
		fmt.Printf("AudioDataAudioTagBody: Varies by format\n")
	}

	return index, nil
}

func (f *Flv) parseAacAudioData(buf []byte, index int) (int, error) {

	var err error
	if f.CurrentTag.AACPacketType == AACPacketTypeAacSequenceHeader {
		index, err = f.parseAudioSpecificConfig(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseAudioSpecificConfig failed, err:%v", err)
		}
	} else if f.CurrentTag.AACPacketType == AACPacketTypeAacRaw {
		index, err = f.parseRawAacFrameData(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseAudioSpecificConfig failed, err:%v", err)
		}

	}
	return index, nil
}

func (f *Flv) parseRawAacFrameData(buf []byte, index int) (int, error) {
	aacFrameLength := f.CurrentTag.Length - index + 7
	byte3 := byte(0x8<<4 + aacFrameLength >> 11)
	byte4 := byte(aacFrameLength >> 3)
	byte5 := byte(aacFrameLength << 5 + 0b00011111)
	byte6 := byte(0b11111100)

	_, _ = aacFile.Write([]byte{0xFF, 0xF1, 0x50,byte3,byte4,byte5,byte6})


	_, _ = aacFile.Write(buf[index:f.CurrentTag.Length])
	fmt.Printf("has Raw AAC frame data but not decode\n")
	return f.CurrentTag.Length, nil
}

func (f *Flv) parseAudioSpecificConfig(buf []byte, index int) (int, error) {

	fmt.Printf("has AudioSpecificConfig but not decode\n")
	return f.CurrentTag.Length, nil
}

func (f *Flv) parseVideoData(buf []byte, index int) (int, error) {

	var err error

	if f.CurrentTag.Filter == FilterPreProcessing {
		index, err = f.parseVideoDataEncryptedBody(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseVideoDataEncryptedBody failed, err:%v", err)
		}
	} else {
		index, err = f.parseVideoDataTagBody(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseVideoDataTagBody failed, err:%v", err)
		}
	}

	return index, nil
}

func (f *Flv) parseVideoDataEncryptedBody(buf []byte, index int) (int, error) {
	return 0, fmt.Errorf("parseVideoDataEncryptedBody error")
}

func (f *Flv) parseVideoDataTagBody(buf []byte, index int) (int, error) {

	var err error

	if f.CurrentTag.FrameType == FrameTypeVideoInfoOrCommandFrame {
		return 0, fmt.Errorf("FrameTypeVideoInfoOrCommandFrame error")
	} else {

		if f.CurrentTag.CodeId == CodecIDSorensonH263 {
			return 0, fmt.Errorf("CodecIDSorensonH263 error")
		}
		if f.CurrentTag.CodeId == CodecIDScreenVideo {
			return 0, fmt.Errorf("CodecIDScreenVideo error")
		}
		if f.CurrentTag.CodeId == CodecIDOn2Vp6 {
			return 0, fmt.Errorf("CodecIDOn2Vp6 error")
		}
		if f.CurrentTag.CodeId == CodecIDOn2Vp6WithAlphaChannel {
			return 0, fmt.Errorf("CodecIDOn2Vp6WithAlphaChannel error")
		}
		if f.CurrentTag.CodeId == CodecIDScreenVideoVersion2 {
			return 0, fmt.Errorf("CodecIDScreenVideoVersion2 error")
		}
		if f.CurrentTag.CodeId == CodecIDAvc {
			index, err = f.parseAvcVideoPacket(buf, index)
			if err != nil {
				return 0, fmt.Errorf("f.parseAvcVideoPacket failed, err:%v", err)
			}
		}
	}

	return index, nil
}

func (f *Flv) parseAvcVideoPacket(buf []byte, index int) (int, error) {

	var err error

	if f.CurrentTag.AVCPacketType == AvcPacketTypeAvcSequenceHeader {
		index, err = f.parseAvcDecoderConfigurationRecord(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseAvcDecoderConfigurationRecord failed, err:%v", err)
		}
	}

	if f.CurrentTag.AVCPacketType == AvcPacketTypeAvcNalu {
		index, err = f.parseOneOrMoreNalus(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseOneOrMoreNalus failed, err:%v", err)
		}
	}

	return index, nil
}

func (f *Flv) parseAvcDecoderConfigurationRecord(buf []byte, index int) (int, error) {
	fmt.Printf("has AvcDecoderConfigurationRecord but not decode\n")
	return f.CurrentTag.Length, nil
}

func (f *Flv) parseOneOrMoreNalus(buf []byte, index int) (int, error) {
	_, _ = h264File.Write([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1})
	_, _ = h264File.Write(buf[index:f.CurrentTag.Length])
	fmt.Printf("has OneOrMoreNalus but not decode\n")
	return f.CurrentTag.Length, nil
}

func (f *Flv) parseScriptData(buf []byte, index int) (int, error) {

	if f.CurrentTag.Filter == FilterPreProcessing { // is Encrypted
		return f.parseScriptDataEncryptedBody(buf, index)
	} else {
		return f.parseScriptDataTagBody(buf, index)
	}
}

func (f *Flv) parseScriptDataEncryptedBody(buf []byte, index int) (int, error) {
	return 0, fmt.Errorf("parseScriptDataEncryptedBody error")
}

func (f *Flv) parseScriptDataTagBody(buf []byte, index int) (int, error) {

	var err error

	index, err = f.parseScriptDataValue(buf, index)
	if err != nil {
		return 0, fmt.Errorf("f.parseScriptDataValue failed, err:%v", err)
	}

	index, err = f.parseScriptDataValue(buf, index)
	if err != nil {
		return 0, fmt.Errorf("f.parseScriptDataValue failed, err:%v", err)
	}

	return index, nil
}

func (f *Flv) parseScriptDataValue(buf []byte, index int) (int, error) {

	if len(buf) < 1 {
		return 0, fmt.Errorf("len(buf) < 1")
	}

	var err error

	valueType := buf[index]
	index += 1

	valueTypeString, ok := ScriptDataValueTypeSet[valueType]
	if !ok {
		return 0, fmt.Errorf("ScriptDataValueTypeSet[valueType] is not ok, valueType:%v", valueType)
	}
	fmt.Printf("Script Data Value Type is %v\n", valueTypeString)

	if valueType == ScriptDataValueTypeNumber {
		if len(buf[index:]) < 8 {
			return 0, fmt.Errorf("ScriptDataValueTypeNumber error: len(buf) < 8")
		}
		doubleValue, err := util.ByteToFloat64(buf[index : index+8])
		if err != nil {
			return 0, fmt.Errorf("util.ByteToFloat64 failed, err:%v ", err)
		}
		fmt.Printf("Script Data Value Number is %f\n", doubleValue)
		index += 8
	}

	if valueType == ScriptDataValueTypeBoolean {
		if len(buf[index:]) < 1 {
			return 0, fmt.Errorf("ScriptDataValueTypeBoolean error: len(buf) < 1")
		}
		booleanValue := util.BytesToUint8ByBigEndian(buf[index])
		fmt.Printf("Script Data Value Boolean is %v\n", booleanValue)

		index += 1
	}

	if valueType == ScriptDataValueTypeString {
		index, err = f.parseScriptDataString(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseScriptDataString failed, err:%v", err)
		}
	}

	if valueType == ScriptDataValueTypeObject {
		return 0, fmt.Errorf("ScriptDataValueTypeObject error")
	}

	if valueType == ScriptDataValueTypeReference {
		return 0, fmt.Errorf("ScriptDataValueTypeReference error")
	}

	if valueType == ScriptDataValueTypeEcmaArray {
		index, err = f.parseScriptDataEcmaArray(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseScriptDataEcmaArray failed, err:%v", err)
		}

	}

	if valueType == ScriptDataValueTypeStrictArray {
		return 0, fmt.Errorf("ScriptDataValueTypeStrictArray error")
	}

	if valueType == ScriptDataValueTypeDate {
		return 0, fmt.Errorf("ScriptDataValueTypeDate error")
	}

	if valueType == ScriptDataValueTypeLongString {
		return 0, fmt.Errorf("ScriptDataValueTypeLongString error")
	}

	return index, nil
}

func (f *Flv) parseScriptDataString(buf []byte, index int) (int, error) {
	if len(buf) < 2 {
		return 0, fmt.Errorf("len(buf) < 2")
	}

	var err error

	stringLength, err := util.BytesToUint16ByBigEndian(buf[index : index+2])
	if err != nil {
		return 0, fmt.Errorf("util.BytesToUint16ByBigEndian failed, err:%v", err)
	}
	//fmt.Printf("stringLength is %v\n", stringLength)
	index += 2

	if len(buf[index:]) < int(stringLength) {
		return 0, fmt.Errorf("parseScriptDataString error")
	}

	stringData := string(buf[index : index+int(stringLength)])
	fmt.Printf("Script Data Value String is %v\n", stringData)

	index += int(stringLength)

	return index, nil
}

func (f *Flv) parseScriptDataEcmaArray(buf []byte, index int) (int, error) {

	if len(buf) < 4 {
		return 0, fmt.Errorf("len(buf) < 4")
	}

	var err error

	ecmaArrayLength, err := util.BytesToUint32ByBigEndian(buf[index : index+4])
	if err != nil {
		return 0, fmt.Errorf("util.BytesToUint32ByBigEndian failed, err:%v", err)
	}
	fmt.Printf("ECMAArrayLength is %v\n", ecmaArrayLength)
	index += 4

	var i int64 = 0
	for ; i < int64(ecmaArrayLength); i++ {
		index, err = f.parseScriptDataString(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseScriptDataString failed, err:%v", err)
		}

		index, err = f.parseScriptDataValue(buf, index)
		if err != nil {
			return 0, fmt.Errorf("f.parseScriptDataValue failed, err:%v", err)
		}
	}

	index, err = f.parseScriptDataObjectEnd(buf, index)
	if err != nil {
		return 0, fmt.Errorf("f.parseScriptDataObjectEnd failed, err:%v", err)
	}

	return index, nil
}

func (f *Flv) parseScriptDataObjectEnd(buf []byte, index int) (int, error) {
	if len(buf) < 3 {
		return 0, fmt.Errorf("len(buf) < 3")
	}

	objectEndMark := buf[index : index+3]
	if util.BytesToUint8ByBigEndian(objectEndMark[0]) != 0 ||
		util.BytesToUint8ByBigEndian(objectEndMark[1]) != 0 ||
		util.BytesToUint8ByBigEndian(objectEndMark[2]) != 9 {
		return 0, fmt.Errorf("objectEndMark is not 0 0 9, objectEndMark:%x", objectEndMark)
	}

	fmt.Println("objectEndMark is 0 0 9")

	index += 3
	return index, nil
}
