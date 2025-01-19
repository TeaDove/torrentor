package torrentor_service

import (
	"context"
	"time"
	"torrentor/schemas"

	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (r *Service) onTorrentComplete(
	ctx context.Context,
	torrentEnt *schemas.TorrentEntityPop,
) error {
	<-torrentEnt.Obj.Complete().On()

	torrentEnt.Completed = true

	//err := r.torrentRepository.TorrentUpdate(
	//	ctx,
	//	torrentEnt.ID,
	//	map[string]any{"completed": true},
	//)
	//if err != nil {
	//	return errors.Wrap(err, "failed in marking torrent complete")
	//}

	zerolog.Ctx(ctx).Info().Dict("torrent", torrentEnt.ZerologDict()).Msg("torrent.ready")

	return nil
}

func (r *Service) onFileCompleteCallback(
	ctx context.Context,
	fileEnt *schemas.FileEntityPop,
) error {
	var err error

	if fileEnt.Mimetype == schemas.MatroskaMimeType {
		err = r.unpackMatroska(ctx, fileEnt)
		if err != nil {
			return errors.Wrap(err, "failed to unpack matroska file")
		}
	}

	//fileEnt.Completed = true

	// TODO file save and update torrent
	//_, err = r.torrentRepository.TorrentInsert(ctx, &fileEnt.Torrent.TorrentEntity)
	//if err != nil {
	//	return errors.Wrap(err, "failed in marking torrent complete")
	//}

	zerolog.Ctx(ctx).
		Debug().
		Interface("file", fileEnt.Name).
		Msg("file.ready")

	return nil
}

func (r *Service) onFileComplete(
	ctx context.Context,
	torrentEnt *schemas.TorrentEntityPop,
	completedCheckPeriod time.Duration,
) error {
	// TODO check if already completed
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
			err := r.onFileCompleteCallback(ctx, &schemas.FileEntityPop{
				FileEntity: torrentEnt.Files[schemas.TrimFirstDir(fileName)],
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
