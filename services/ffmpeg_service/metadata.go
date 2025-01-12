package ffmpeg_service

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Metadata struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	Index         int       `json:"index"`
	CodecName     string    `json:"codec_name"`
	CodecLongName string    `json:"codec_long_name"`
	CodecType     CodecType `json:"codec_type"`
	Tags          Tag       `json:"tags"`
}

type CodecType string

const (
	CodecTypeAudio    = "audio"
	CodecTypeVideo    = "video"
	CodecTypeSubtitle = "subtitle"
)

type Tag struct {
	Language string `json:"language,omitempty"`
	Title    string `json:"title,omitempty"`
}

func (r *Service) ExportMetadata(ctx context.Context, filePath string) (Metadata, error) {
	metadataRaw, err := ffmpeg.Probe(filePath)
	if err != nil {
		return Metadata{}, errors.Wrap(err, "failed to fetch ffmpeg metadata")
	}

	var metadata Metadata

	err = json.Unmarshal([]byte(metadataRaw), &metadata)
	if err != nil {
		return Metadata{}, errors.Wrap(err, "failed to unmarshal ffmpeg metadata")
	}

	zerolog.Ctx(ctx).
		Info().
		Interface("metadata", metadata).
		Msg("ffmpeg.metadata.exported")

	return metadata, nil
}
