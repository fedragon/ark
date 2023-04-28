package header

const (
	dateTimeOriginal   = 0x9003 // Date/time when original image was taken
	offsetTimeOriginal = 0x9011 // Name of dateTimeOriginal's timezone (e.g. Europe/Amsterdam)
	timeZoneOffset     = 0x882a // Offset of dateTimeOriginal's timezone from GMT, in hours
)
