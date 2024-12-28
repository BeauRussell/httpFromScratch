package framePackaging

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/net/http2/hpack"
)

type FrameType uint8

const (
	maxFrameSize                = 16384
	FrameData         FrameType = 0x00
	FrameHeaders      FrameType = 0x01
	FramePriority     FrameType = 0x02
	FrameRSTStream    FrameType = 0x03
	FrameSettings     FrameType = 0x04
	FramePushPromise  FrameType = 0x05
	FramePing         FrameType = 0x06
	FrameGoAway       FrameType = 0x07
	FrameWindowUpdate FrameType = 0x08
	FrameContinuation FrameType = 0x09
	FrameAltSVC       FrameType = 0x0a
	FrameOrigin       FrameType = 0x0c
)

type Frame struct {
	Length   uint32
	Type     FrameType
	Flags    uint8
	StreamID uint32
	Payload  []byte
}

func (f *Frame) ParseFrame(reader io.Reader) error {
	err := f.parseFrameHeader(reader)
	if err != nil {
		fmt.Println("Failed to parse frame headers:", err)
		return err
	}

	if f.Length > 0 {
		f.Payload = make([]byte, f.Length)
		if _, err := io.ReadFull(reader, f.Payload); err != nil {
			fmt.Println("Failed to read frame payload:", err)
			return err
		}
	}

	return nil
}

func (f *Frame) parseFrameHeader(reader io.Reader) error {
	header := make([]byte, 9)
	if _, err := io.ReadFull(reader, header); err != nil {
		fmt.Println("Failed to read into headers:", err)
		return err
	}

	f.Length = uint32(header[0])<<16 | uint32(header[1])<<8 | uint32(header[2])
	f.Type = FrameType(header[3])
	f.Flags = header[4]
	f.StreamID = binary.BigEndian.Uint32(header[5:]) & 0x7FFFFFFF // Mask out the reserved bit

	if f.Length > maxFrameSize {
		return fmt.Errorf("Frame length exceed maximum allowed size: %d", f.Length)
	}
	return nil
}

func (f *Frame) BuildHeadersFrame(headers map[string]string, streamID uint32) ([]byte, error) {
	f.Type = FrameHeaders
	f.StreamID = streamID
	f.Flags = 0x04 // END_HEADERS

	var headerBuf bytes.Buffer
	encoder := hpack.NewEncoder(&headerBuf)
	for key, value := range headers {
		if err := encoder.WriteField(hpack.HeaderField{Name: key, Value: value}); err != nil {
			return nil, fmt.Errorf("failed to encode header field: %w", err)
		}
	}
	f.Payload = headerBuf.Bytes()
	f.Length = uint32(len(f.Payload))

	buf := bytes.Buffer{}

	lenBytes, err := convertLengthToBytes(f.Length)
	if err != nil {
		return nil, err
	}
	buf.Write(lenBytes)
	buf.WriteByte(byte(f.Type))
	buf.WriteByte(f.Flags)
	binary.Write(&buf, binary.BigEndian, streamID&0x7FFFFFFF)

	buf.Write(f.Payload)

	return buf.Bytes(), nil
}

func (f *Frame) BuildDataFrame(data string, stream uint32) ([]byte, error) {
	f.Type = FrameData
	f.Length = uint32(len(data))
	f.Flags = 0x01 // END_STREAM
	f.StreamID = stream
	f.Payload = []byte(data)

	buf := bytes.Buffer{}

	lenBytes, err := convertLengthToBytes(f.Length)
	if err != nil {
		return nil, err
	}

	buf.Write(lenBytes)
	buf.WriteByte(0x00)
	buf.WriteByte(f.Flags)
	binary.Write(&buf, binary.BigEndian, stream&0x7FFFFFFF)

	buf.Write(f.Payload)

	return buf.Bytes(), nil
}

func convertLengthToBytes(length uint32) ([]byte, error) {
	if length > 0xFFFFFF {
		return nil, fmt.Errorf("length to convert exceeds 3-byte limit: %d", length)
	}

	return []byte{
		byte((length >> 16) & 0xFF),
		byte((length >> 8) & 0xFF),
		byte(length & 0xFF),
	}, nil
}
