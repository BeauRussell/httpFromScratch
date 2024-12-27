package framePackaging

import (
	"encoding/binary"
	"fmt"
	"io"
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
