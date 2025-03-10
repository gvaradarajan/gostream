// Package codec defines the encoder and factory interfaces for encoding video frames.
package codec

import (
	"image"

	"github.com/edaniels/golog"
)

// DefaultKeyFrameInterval is the default interval chosen
// in order to produce high enough quality results at a low
// latency.
const DefaultKeyFrameInterval = 30

// An Encoder is anything that can encode images into bytes. This means that
// the encoder must follow some type of format dictated by a type (see EncoderFactory.MimeType).
// An encoder that produces bytes of different encoding formats per call is invalid.
type Encoder interface {
	Encode(img image.Image) ([]byte, error)
}

// An EncoderFactory produces Encoders and provides information about the underlying encoder itself.
type EncoderFactory interface {
	New(height, width, keyFrameInterval int, logger golog.Logger) (Encoder, error)
	MIMEType() string
}
