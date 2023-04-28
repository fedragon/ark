package exif

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

const (
	OrfHeaderSize     = 16
	IntelByteOrder    = 0x4949
	MotorolaByteOrder = 0x4D4D

	TiffMagicNumberBigEndian    = 0x002A
	TiffMagicNumberLittleEndian = 0x2A00
	OrfMagicNumberBigEndian     = 0x4F52
	OrfMagicNumberLittleEndian  = 0x524F
)

// parseFromOrf parses CreatedAt from the EXIF data of an ORF file.
func parseFromOrf(r io.ReadSeeker) (time.Time, error) {
	header := make([]byte, OrfHeaderSize)
	_, err := r.Read(header)
	if err != nil {
		return time.Time{}, err
	}

	byteOrder, err := parseEndianness(header[0:2])
	if err != nil {
		return time.Time{}, err
	}

	if err := parseMagicNumber(byteOrder, header[2:4]); err != nil {
		return time.Time{}, err
	}

	return time.Time{}, nil
}

func parseEndianness(bytes []byte) (binary.ByteOrder, error) {
	// Note: the value of these 2 bytes is endianess-independent, so I can use any byte order to read them.
	endianness := binary.LittleEndian.Uint16(bytes)
	switch endianness {
	case IntelByteOrder:
		return binary.LittleEndian, nil
	case MotorolaByteOrder:
		return binary.BigEndian, nil
	default:
		return nil, fmt.Errorf("unknown endianess: 0x%X", endianness)
	}
}

func parseMagicNumber(byteOrder binary.ByteOrder, bytes []byte) error {
	magicNumber := byteOrder.Uint16(bytes)
	if magicNumber != TiffMagicNumberBigEndian &&
		magicNumber != TiffMagicNumberLittleEndian &&
		magicNumber != OrfMagicNumberBigEndian &&
		magicNumber != OrfMagicNumberLittleEndian {
		return fmt.Errorf("unknown magic number: 0x%X", magicNumber)
	}
	return nil
}
