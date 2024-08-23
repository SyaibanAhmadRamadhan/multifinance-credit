package primitive

type MimeType string

const (
	MimeTypePng  MimeType = "image/png"
	MimeTypeJpeg MimeType = "image/jpeg"
	MimeTypeGif  MimeType = "image/gif"
)

var AllowedMimeTypes = []MimeType{
	MimeTypePng,
	MimeTypeJpeg,
	MimeTypeGif,
}

var MapMimeTypeExtensions = map[MimeType]string{
	MimeTypePng: ".png",
	MimeTypeGif: ".gif",
}

func (v MimeType) IsValid() bool {
	for _, existing := range AllowedMimeTypes {
		if existing == v {
			return true
		}
	}
	return false
}
