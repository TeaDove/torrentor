package ffmpeg_service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Metadata struct {
	Streams   []Stream          `json:"streams"`
	StreamMap map[string]Stream `json:"-"`
}

func (r *Metadata) AudioStreamsAsStrings() []string {
	v := make([]string, 0, len(r.Streams))
	for _, stream := range r.Streams {
		if stream.CodecType == CodecTypeAudio {
			v = append(v, stream.String())
		}
	}

	return v
}

type Stream struct {
	Index         int       `json:"index"`
	CodecName     string    `json:"codec_name"`
	CodecLongName string    `json:"codec_long_name"`
	CodecType     CodecType `json:"codec_type"`
	Tags          Tag       `json:"tags"`
}

func (r *Stream) StreamFile(suffix string) string {
	return fmt.Sprintf("%s%s", r, suffix)
}

func (r *Stream) String() string {
	fields := []string{strconv.Itoa(r.Index)}
	if r.Tags.Title != "" {
		fields = append(fields, r.Tags.Title)
	}
	if r.Tags.Language != "" {
		fields = append(fields, r.Tags.Language)
	}

	return strings.Join(fields, "-")
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

func (r *Service) ExportMetadata(_ context.Context, filePath string) (Metadata, error) {
	metadataRaw, err := ffmpeg.Probe(filePath)
	if err != nil {
		return Metadata{}, errors.Wrap(err, "failed to fetch ffmpeg metadata")
	}

	var metadata Metadata

	err = json.Unmarshal([]byte(metadataRaw), &metadata)
	if err != nil {
		return Metadata{}, errors.Wrap(err, "failed to unmarshal ffmpeg metadata")
	}

	metadata.StreamMap = make(map[string]Stream, len(metadata.Streams))
	for _, stream := range metadata.Streams {
		metadata.StreamMap[stream.String()] = stream
	}

	return metadata, nil
}
