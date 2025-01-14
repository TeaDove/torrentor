package torrentor_service

import (
	"context"
	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"time"
	"torrentor/schemas"
)

func (r *Service) onTorrentComplete(
	ctx context.Context,
	torrentEnt *schemas.TorrentEntityPop,
) error {
	<-torrentEnt.Obj.Complete().On()

	zerolog.Ctx(ctx).Info().Dict("torrent", torrentEnt.ZerologDict()).Msg("torrent.ready")

	torrentEnt.Completed = true

	_, err := r.torrentRepository.TorrentUpsert(ctx, &torrentEnt.TorrentEntity)
	if err != nil {
		return errors.Wrap(err, "failed in marking torrent complete")
	}

	return nil
}

func (r *Service) onFileCompleteCallback(
	ctx context.Context,
	fileEnt *schemas.FileEntityPop,
) error {
	if fileEnt.Mimetype != schemas.MatroskaMimeType {
		return nil
	}

	err := r.unpackMatroska(ctx, fileEnt)
	if err != nil {
		return errors.Wrap(err, "failed to unpack matroska file")
	}

	return nil
}

func (r *Service) onFileComplete(
	ctx context.Context,
	torrentEnt *schemas.TorrentEntityPop,
	completedCheckPeriod time.Duration,
) error {
	incompleteFiles := map[string]*torrent.File{}
	for _, file := range torrentEnt.Obj.Files() {
		incompleteFiles[file.Path()] = file
	}
	completed := make([]string, 0, len(incompleteFiles))

	for {
		for _, file := range incompleteFiles {
			if file.Length() == file.BytesCompleted() {
				completed = append(completed, file.Path())
			}
		}
		for _, fileName := range completed {
			zerolog.Ctx(ctx).
				Debug().
				Str("file", fileName).
				Msg("file.ready")

			err := r.onFileCompleteCallback(ctx, &schemas.FileEntityPop{
				FileEntity: torrentEnt.Files[fileName],
				Obj:        incompleteFiles[fileName],
				Torrent:    torrentEnt,
			})
			if err != nil {
				return errors.Wrap(err, "failed to unpack matroska")
			}

			delete(incompleteFiles, fileName)
		}
		completed = make([]string, 0, len(incompleteFiles))
		if len(incompleteFiles) == 0 {
			break
		}

		time.Sleep(completedCheckPeriod)
	}

	err := r.onTorrentComplete(ctx, torrentEnt)
	if err != nil {
		return errors.Wrap(err, "failed to mark complete")
	}

	return nil
}
