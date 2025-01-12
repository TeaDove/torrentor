package ffmpeg_service

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type metadata struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Index         int    `json:"index"`
	CodecName     string `json:"codec_name"`
	CodecLongName string `json:"codec_long_name"`
	CodecType     string `json:"codec_type"`
	Tags          tag    `json:"tags"`
}

type codecType string

const (
	codecTypeAudio    = "audio"
	codecTypeVideo    = "video"
	codecTypeSubtitle = "subtitle"
)

type tag struct {
	Language string `json:"language,omitempty"`
	Title    string `json:"title,omitempty"`
}

func (r *Service) exportMetadata(ctx context.Context, filePath string) (metadata, error) {
	metadataRaw, err := ffmpeg.Probe(filePath)
	if err != nil {
		return metadata{}, errors.Wrap(err, "failed to fetch ffmpeg metadata")
	}

	var metadataObj metadata
	err = json.Unmarshal([]byte(metadataRaw), &metadataObj)
	if err != nil {
		return metadata{}, errors.Wrap(err, "failed to unmarshal ffmpeg metadata")
	}

	zerolog.Ctx(ctx).
		Info().
		Interface("metadata", metadataObj).
		Msg("ffmpeg.metadata.exported")

	return metadataObj, nil
}
