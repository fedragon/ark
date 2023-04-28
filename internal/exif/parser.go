package exif

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	exifcommon "github.com/dsoprea/go-exif/v2/common"
	heicexif "github.com/dsoprea/go-heic-exif-extractor"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure"
	tiffstructure "github.com/dsoprea/go-tiff-image-structure"
	riimage "github.com/dsoprea/go-utility/image"
)

var parsers = map[string]riimage.MediaParser{
	".jpg":  jpegstructure.NewJpegMediaParser(),
	".jpeg": jpegstructure.NewJpegMediaParser(),
	".cr2":  tiffstructure.NewTiffMediaParser(),
	".tiff": tiffstructure.NewTiffMediaParser(),
	".heic": heicexif.NewHeicExifMediaParser(),
}

func ParseCreatedAt(path string) (time.Time, error) {
	ext := strings.ToLower(filepath.Ext(path))

	var createdAt time.Time
	var err error
	var done bool

	if parser, ok := parsers[ext]; ok {
		createdAt, err, done = parse(parser, path)
	}

	if err != nil {
		return time.Time{}, err
	}

	if done {
		return createdAt, nil
	}

	fmt.Printf("Unknown extension: '%s'. Using file modification time instead.\n", ext)
	stat, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	return stat.ModTime(), nil
}

func parse(parser riimage.MediaParser, path string) (time.Time, error, bool) {
	ctx, err := parser.ParseFile(path)
	if err != nil {
		return time.Time{}, err, false
	}
	ifd, _, err := ctx.Exif()
	if err != nil {
		return time.Time{}, err, false
	}
	exif, err := ifd.ChildWithIfdPath(exifcommon.IfdPathStandardExif)
	if err != nil {
		return time.Time{}, err, false
	}
	tags, err := exif.FindTagWithId(0x9003)
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
