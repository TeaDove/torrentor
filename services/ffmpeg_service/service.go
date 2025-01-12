package ffmpeg_service

import "context"

type Service struct{}

func NewService(_ context.Context) (*Service, error) {
	return &Service{}, nil
}
