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
	torrentEnt *schemas.TorrentEntity,
) error {
	<-torrentEnt.Obj.Complete().On()

	torrentEnt.Completed = true

	zerolog.Ctx(ctx).
		Info().
		Dict("torrent", torrentEnt.ZerologDict()).
		Msg("torrent.download.completed")

	return nil
}

func (r *Service) onFileCompleteCallback(
	ctx context.Context,
	fileEnt *schemas.FileEntity,
	unpack bool,
) error {
	if unpack {
		err := r.UnpackIfNeeded(ctx, fileEnt)
		if err != nil {
			return errors.Wrap(err, "failed to unpack")
		}
	}

	fileEnt.Completed = true

	zerolog.Ctx(ctx).
		Debug().
		Dict("file", fileEnt.ZerologDict()).
		Msg("file.ready")

	return nil
}

func (r *Service) onFileComplete(
	ctx context.Context,
	torrentEnt *schemas.TorrentEntity,
	completedCheckPeriod time.Duration,
) error {
	allFiles := torrentEnt.FlatFiles()
	if len(allFiles) == 0 {
		panic("torrent has no files")
	}
	firstFilePath := allFiles[0].Obj.Path()

	// TODO check if already completed
	incompleteFiles := map[string]*torrent.File{}
	for _, file := range torrentEnt.Obj.Files() {
		incompleteFiles[file.Path()] = file
	}

	completedFilePaths := make([]string, 0, len(incompleteFiles))

	for {
		for _, file := range incompleteFiles {
			if file.Length() == file.BytesCompleted() {
				completedFilePaths = append(completedFilePaths, file.Path())
			}
		}

		for _, filePath := range completedFilePaths {
			unpack := false
			if filePath == firstFilePath {
				unpack = true
			}

			err := r.onFileCompleteCallback(
				ctx,
				torrentEnt.FilePathMap[schemas.TrimFirstDir(filePath)],
				unpack,
			)
			if err != nil {
				return errors.Wrap(err, "failed to unpack matroska")
			}

			delete(incompleteFiles, filePath)
		}

		completedFilePaths = make([]string, 0, len(incompleteFiles))

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
