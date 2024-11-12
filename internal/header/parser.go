package header

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	common "github.com/dsoprea/go-exif/v2/common"
	heic "github.com/dsoprea/go-heic-exif-extractor"
	jpeg "github.com/dsoprea/go-jpeg-image-structure"
	image "github.com/dsoprea/go-utility/image"
	"github.com/fedragon/tiff-parser/tiff"
)

var parsers = map[string]image.MediaParser{
	".jpg":  jpeg.NewJpegMediaParser(),
	".jpeg": jpeg.NewJpegMediaParser(),
	".heic": heic.NewHeicExifMediaParser(),
}

var tiffs = map[string]struct{}{
	".cr2":  {},
	".orf":  {},
	".tiff": {},
}

func ParseCreatedAt(path string) (time.Time, error) {
	ext := strings.ToLower(filepath.Ext(path))

	var createdAt time.Time
	var err error
	var done bool

	if parser, ok := parsers[ext]; ok {
		createdAt, err, done = parse(parser, path)
	}

	if err != nil && err.Error() != "no exif data" {
		return time.Time{}, err
	}

	if done {
		return createdAt, nil
	}

	if _, ok := tiffs[ext]; ok {
		var reader io.ReadSeekCloser
		reader, err = os.Open(path)
		if err != nil {
			return time.Time{}, err
		}
		defer reader.Close()

		parser, err := tiff.NewParser(reader)
		if err != nil {
			return time.Time{}, err
		}

		entries, err := parser.Parse(tiff.DateTimeOriginal)
		if err != nil {
			return time.Time{}, err
		}

		if en, ok := entries[tiff.DateTimeOriginal]; ok {
			switch en.DataType {
			case tiff.DataType_String:
				return time.Parse("2006:01:02 15:04:05", *en.Value.String)
			}
		}
	}

	return time.Time{}, notFound(ext)
}

func parse(parser image.MediaParser, path string) (time.Time, error, bool) {
	ctx, err := parser.ParseFile(path)
	if err != nil {
		return time.Time{}, err, false
	}
	ifd, _, err := ctx.Exif()
	if err != nil {
		return time.Time{}, err, false
	}
	exif, err := ifd.ChildWithIfdPath(common.IfdPathStandardExif)
	if err != nil {
		return time.Time{}, err, false
	}
	tags, err := exif.FindTagWithId(0x9003) // dateTimeOriginal
	if err != nil {
		return time.Time{}, err, false
	}

	for _, tag := range tags {
		value, err := tag.Value()
		if err != nil {
			return time.Time{}, err, false
		}

		parsed, err := time.Parse("2006:01:02 15:04:05", value.(string))
		return parsed, err, true
	}

	return time.Time{}, nil, false
}
