package header

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	IntelByteOrder    = 0x4949
	MotorolaByteOrder = 0x4D4D

	TiffMagicNumberBigEndian    = 0x002A
	TiffMagicNumberLittleEndian = 0x2A00
	OrfMagicNumberBigEndian     = 0x4F52
	OrfMagicNumberLittleEndian  = 0x524F

	ExifOffsetId = 0x8769
)

// wantedIDs represents a set of Tag IDs
type wantedIDs struct {
	wanted map[uint16]struct{}
	max    uint16
}

func newWantedIDs(ids ...uint16) *wantedIDs {
	we := wantedIDs{
		wanted: map[uint16]struct{}{},
	}
	for _, id := range ids {
		we.Put(id)
	}

	return &we
}

func (we *wantedIDs) Put(id uint16) {
	we.wanted[id] = struct{}{}
	if id > we.max {
		we.max = id
	}
}

func (we *wantedIDs) Contains(id uint16) bool {
	_, ok := we.wanted[id]
	return ok
}

func (we *wantedIDs) Max() uint16 {
	return we.max
}

// Tag represents an IFD entry (aka Tag)
type Tag struct {
	ID       uint16
	DataType uint16
	Length   uint32
	Value    uint32 // value of the entry or byte offset to read the value from, depending on the DataType and Length
}

// parseFromTiff parses CreatedAt from the EXIF data of a TIFF file. Several vendor-specific formats conform to the TIFF header structure (e.g. CR2, ORF).
func parseFromTiff(r io.ReadSeeker) (time.Time, bool, error) {
	header := make([]byte, 8) // only read the first 3 fields of the TIFF header
	_, err := io.ReadFull(r, header)
	if err != nil {
		return time.Time{}, false, err
	}

	byteOrder, err := readEndianness(header[0:2])
	if err != nil {
		return time.Time{}, false, err
	}

	if err := validateMagicNumber(byteOrder, header[2:4]); err != nil {
		return time.Time{}, false, err
	}

	offsetToFirstIfd := int64(byteOrder.Uint32(header[4:8]))
	if _, err := r.Seek(offsetToFirstIfd, io.SeekStart); err != nil {
		return time.Time{}, false, err
	}

	dt, err := readCreatedAt(byteOrder, r, offsetToFirstIfd)
	if err != nil {
		return time.Time{}, false, err
	}

	return dt, true, nil
}

// readEndianness reads and returns the endianness of the metadata.
func readEndianness(buffer []byte) (binary.ByteOrder, error) {
	// Note: the value of these 2 bytes is endianness-independent, so I can use any byte order to read them.
	value := binary.LittleEndian.Uint16(buffer)
	switch value {
	case IntelByteOrder:
		return binary.LittleEndian, nil
	case MotorolaByteOrder:
		return binary.BigEndian, nil
	default:
		return nil, fmt.Errorf("unknown endianness: 0x%X", value)
	}
}

// validateMagicNumber validates the file type by checking that it conforms to one of the expected values
func validateMagicNumber(byteOrder binary.ByteOrder, buffer []byte) error {
	magicNumber := byteOrder.Uint16(buffer)
	if magicNumber != TiffMagicNumberBigEndian &&
		magicNumber != TiffMagicNumberLittleEndian &&
		magicNumber != OrfMagicNumberBigEndian &&
		magicNumber != OrfMagicNumberLittleEndian {
		return fmt.Errorf("unknown magic number: 0x%X", magicNumber)
	}
	return nil
}

// readCreatedAt reads and returns the original date/time from the EXIF subdirectory of IFD#0.
func readCreatedAt(byteOrder binary.ByteOrder, r io.ReadSeeker, offset int64) (time.Time, error) {
	entries, err := collectIFDEntries(byteOrder, r, offset, newWantedIDs(ExifOffsetId))
	if err != nil {
		return time.Time{}, err
	}

	exifOffset, ok := entries[ExifOffsetId]
	if !ok {
		return time.Time{}, errors.New("no exif data")
	}

	entries, err = collectIFDEntries(byteOrder, r, int64(exifOffset.Value), newWantedIDs(dateTimeOriginal, offsetTimeOriginal, timeZoneOffset))
	if err != nil {
		return time.Time{}, err
	}

	dateTimeOriginal, ok := entries[dateTimeOriginal]
	if !ok {
		return time.Time{}, errors.New("dateTimeOriginal not found")
	}

	s, err := readString(dateTimeOriginal, r)
	if err != nil {
		return time.Time{}, err
	}

	dateTime, err := time.Parse("2006:01:02 15:04:05", s)
	if err != nil {
		return time.Time{}, err
	}

	return dateTime, nil
}

// collectIFDEntries collects a set of IFD entries from an IFD.
// To save memory and time (an IFD may contain tens of thousands of entries), it returns as soon as:
// - all entries have been collected, or
// - it has scanned the maximum ID among the desired ones (tags are written according to the natural ordering of their
// ID value: no point in looking further).
func collectIFDEntries(byteOrder binary.ByteOrder, r io.ReadSeeker, offset int64, wanted *wantedIDs) (map[uint16]Tag, error) {
	var entries = make(map[uint16]Tag)

	buffer := make([]byte, 2)
	_, err := r.Read(buffer)
	if err != nil {
		return nil, err
	}
	numEntries := int64(byteOrder.Uint16(buffer))

	for i := int64(0); i < numEntries; i++ {
		buffer := make([]byte, 12)
		if _, err := r.Seek(offset+2+i*12, io.SeekStart); err != nil {
			return nil, err
		}
		_, err := io.ReadFull(r, buffer)
		if err != nil {
			return nil, err
		}

		id := byteOrder.Uint16(buffer[:2])
		if wanted.Contains(id) {
			entries[id] = Tag{
				ID:       id,
				DataType: byteOrder.Uint16(buffer[2:4]),
				Length:   byteOrder.Uint32(buffer[4:8]),
				Value:    byteOrder.Uint32(buffer[8:12]),
			}
		}

		// No point in scanning the IFD further: if we've already found all desired IDs, we're done; if not, we're not going to find them further anyway
		if id >= wanted.Max() {
			break
		}
	}

	return entries, nil
}

// readString reads and returns a string from an IFD entry, trimming its NUL-byte terminator.
func readString(entry Tag, r io.ReadSeeker) (string, error) {
	if _, err := r.Seek(int64(entry.Value), io.SeekStart); err != nil {
		return "", err
	}
	res := make([]byte, entry.Length)
	if _, err := r.Read(res); err != nil {
		return "", err
	}

	return string(bytes.TrimSuffix(res, []byte{0x0})), nil
}
