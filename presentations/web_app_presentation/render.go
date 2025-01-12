package web_app_presentation

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func viewError(c fiber.Ctx, err error) error {
	zerolog.Ctx(c.Context()).
		Error().
		Stack().Err(err).
		Msg("view.error")
	c.Status(fiber.StatusInternalServerError)

	return c.Render("error", fiber.Map{"Error": errors.Wrap(err, "failed to parse id")})
}

func IndexForm(c fiber.Ctx) error {
	return c.Render("index", fiber.Map{"IP": c.IP()})
}

func getParamsUUID(c fiber.Ctx, key string) (uuid.UUID, error) {
	idStr := c.Params(key)

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, viewError(c, errors.Wrap(err, "failed to parse id"))
	}

	return id, nil
}

func (r *Presentation) TorrentForm(c fiber.Ctx) error {
	id, err := getParamsUUID(c, "id")
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to parse id"))
	}

	torrent, err := r.torrentorService.GetTorrentMetadataByID(c.Context(), id)
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to get torrent metadata"))
	}

	return c.Render("torrent",
		fiber.Map{
			"TorrentName":  torrent.Name,
			"TorrentFiles": torrent.FlatFiles(),
			"TorrentId":    id,
		},
	)
}

func (r *Presentation) FileForm(c fiber.Ctx) error {
	torrentID, err := getParamsUUID(c, "id")
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to parse torrent"))
	}

	filePath := c.Query("path")
	if filePath == "" {
		return viewError(c, errors.New("no file path specified"))
	}

	file, err := r.torrentorService.GetFileWithContent(c.Context(), torrentID, filePath)
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to get file content"))
	}

	if file.IsDir {
		return viewError(c, errors.New("file is a directory"))
	}

	mimeType := file.Mimetype
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	c.Set("Content-Type", mimeType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))

	return c.SendStream(file.OSFile, int(file.Size))
}

func (r *Presentation) WatchForm(c fiber.Ctx) error {
	torrentID, err := getParamsUUID(c, "id")
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to parse torrent"))
	}

	filePath := c.Query("path")
	if filePath == "" {
		return viewError(c, errors.New("no file path specified"))
	}

	fileMeta, err := r.torrentorService.GetFile(c.Context(), torrentID, filePath)
	if err != nil {
		return viewError(c, errors.Wrap(err, "failed to get file content"))
	}

	return c.Render("watch",
		fiber.Map{
			"TorrentID": torrentID,
			"Path":      fileMeta.Path,
			"Mimetype":  fileMeta.Mimetype,
		},
	)
}
