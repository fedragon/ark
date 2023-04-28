# TODO

## Extract EXIF data from images

- [ ] ORF
- [x] CR2, TIFF
- [x] JPG, JPEG
- [x] HEIC
- [x] Fallback to Modified-Date if no EXIF data is found
- [ ] Consider writing my own parser for all the above formats, because the library I'm currently using parses the whole file (which is likely to be tens of MBs)
