package header

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseEndianness(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
		order binary.ByteOrder
		err   bool
	}{
		{
			name:  "IntelByteOrder",
			input: []byte{0x49, 0x49},
			order: binary.LittleEndian,
			err:   false,
		},
		{
			name:  "MotorolaByteOrder",
			input: []byte{0x4D, 0x4D},
			order: binary.BigEndian,
			err:   false,
		},
		{
			name:  "UnknownByteOrder",
			input: []byte{0x34, 0x4D},
			order: nil,
			err:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order, err := readEndianness(tc.input)
			if tc.err && err == nil {
				t.Error("expected error, but got none")
			}
			if !tc.err && err != nil {
				t.Error(err)
			}
			if order != tc.order {
				t.Errorf("expected order %v, but got %v", tc.order, order)
			}
		})
	}
}

func Test_ParseMagicNumber(t *testing.T) {
	testCases := []struct {
		name      string
		byteOrder binary.ByteOrder
		input     []byte
		err       bool
	}{
		{
			name:      "TiffMagicNumberBigEndian",
			byteOrder: binary.BigEndian,
			input:     []byte{0x00, 0x2A},
			err:       false,
		},
		{
			name:      "TiffMagicNumberLittleEndian",
			byteOrder: binary.LittleEndian,
			input:     []byte{0x2A, 0x00},
			err:       false,
		},
		{
			name:      "OrfMagicNumberBigEndian",
			byteOrder: binary.BigEndian,
			input:     []byte{0x4F, 0x52},
			err:       false,
		},
		{
			name:      "OrfMagicNumberLittleEndian",
			byteOrder: binary.LittleEndian,
			input:     []byte{0x52, 0x4F},
			err:       false,
		},
		{
			name:      "UnknownMagicNumber",
			byteOrder: binary.BigEndian,
			input:     []byte{0x34, 0x12},
			err:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateMagicNumber(tc.byteOrder, tc.input)
			if tc.err && err == nil {
				t.Error("expected error, but got none")
			} else if !tc.err && err != nil {
				t.Error(err)
			}
		})
	}
}

func TestParseCreatedAt_CR2(t *testing.T) {
	r, err := os.Open("../../test/data/a/image.cr2")
	assert.NoError(t, err)

	_, _, err = parseFromTiff(r)
	if err != nil {
		assert.NoError(t, err)
	}
}

func TestParseCreatedAt_ORF(t *testing.T) {
	r, err := os.Open("../../test/data/a/image.orf")
	assert.NoError(t, err)

	_, _, err = parseFromTiff(r)
	if err != nil {
		assert.NoError(t, err)
	}
}
