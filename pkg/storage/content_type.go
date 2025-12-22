package storage

import (
	"strings"

	"github.com/jorgeAM/go-template/pkg/errors"
)

var (
	ErrInvalidContentType                 = errors.Define("content_type.internal_content_type")
	ErrContentTypeDoesNotHaveTypeAssigned = errors.Define("content_type.does_not_have_type_assigned")
)

type ContentType string

const (
	IMAGE_JPEG      ContentType = "image/jpeg"
	IMAGE_PNG       ContentType = "image/png"
	IMAGE_HEIC      ContentType = "image/heic"
	IMAGE_HEIF      ContentType = "image/heif"
	IMAGE_WEBP      ContentType = "image/webp"
	TEXT_CSV        ContentType = "text/csv"
	APPLICATION_PDF ContentType = "application/pdf"
)

var allowedContentType = map[string]ContentType{
	IMAGE_JPEG.String():      IMAGE_JPEG,
	IMAGE_PNG.String():       IMAGE_PNG,
	IMAGE_HEIC.String():      IMAGE_HEIC,
	IMAGE_HEIF.String():      IMAGE_HEIF,
	IMAGE_WEBP.String():      IMAGE_WEBP,
	TEXT_CSV.String():        TEXT_CSV,
	APPLICATION_PDF.String(): APPLICATION_PDF,
}

func NewContentType(c string) (ContentType, error) {
	if contentType, ok := allowedContentType[strings.ToLower(c)]; ok {
		return contentType, nil
	}

	return "", errors.New(
		ErrInvalidContentType,
		"invalid content type",
		errors.WithMetadata("content_type", c),
	)
}

func (c ContentType) String() string {
	return string(c)
}

var typeByContentType = map[ContentType]string{
	IMAGE_JPEG:      "jpeg",
	IMAGE_PNG:       "png",
	IMAGE_HEIC:      "heic",
	IMAGE_HEIF:      "heif",
	IMAGE_WEBP:      "webp",
	TEXT_CSV:        "csv",
	APPLICATION_PDF: "pdf",
}

func (c ContentType) GetType() (string, error) {
	if t, ok := typeByContentType[c]; ok {
		return t, nil
	}

	return "", errors.New(
		ErrContentTypeDoesNotHaveTypeAssigned,
		"conntent type does not have type assigned",
		errors.WithMetadata("content_type", c),
	)
}
