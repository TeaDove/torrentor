package web_app_presentation

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
	"path/filepath"
)

func IndexForm(c fiber.Ctx) error {
	return c.Render("index", fiber.Map{"IP": c.IP()})
}

func getParamsInfoHash(c fiber.Ctx) (metainfo.Hash, error) {
	infoHashRaw := c.Params("infohash")

	hash := metainfo.Hash{}
	err := hash.FromHexString(infoHashRaw)
	if err != nil {
		return metainfo.Hash{}, errors.Wrap(err, "failed to parse infohash")
	}

	return hash, nil
}

func (r *Presentation) TorrentForm(c fiber.Ctx) error {
	infoHash, err := getParamsInfoHash(c)
	if err != nil {
		return errors.Wrap(err, "failed to parse infohash")
	}

	torrent, err := r.torrentorService.GetTorrentByInfoHash(c.Context(), infoHash)
	if err != nil {
		return errors.Wrap(err, "failed to get torrent metadata")
	}

	return c.Render("torrent",
		fiber.Map{
			"TorrentName":     torrent.Name,
			"TorrentFiles":    torrent.FlatFiles(),
			"TorrentInfoHash": infoHash,
		},
	)
}

func (r *Presentation) FileForm(c fiber.Ctx) error {
	torrentInfoHash, err := getParamsInfoHash(c)
	if err != nil {
		return errors.Wrap(err, "failed to parse torrent")
	}

	filePath := c.Query("path")
	if filePath == "" {
		return errors.New("no file path specified")
	}

	file, err := r.torrentorService.GetFileByInfoHashAndPath(c.Context(), torrentInfoHash, filePath)
	if err != nil {
		return errors.Wrap(err, "failed to get file content")
	}

	//mimeType := file.Mimetype
	//if mimeType == "" {
	//	mimeType = "application/octet-stream"
	//}
	//
	//c.Set("Content-Type", mimeType)
	//c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))

	return c.SendFile(file.Location())
}

func (r *Presentation) HLSForm(c fiber.Ctx) error {
	torrentInfoHash, err := getParamsInfoHash(c)
	if err != nil {
		return errors.Wrap(err, "failed to parse torrent")
	}

	fileHash := c.Params("filehash")
	if fileHash == "" {
		return errors.New("no file path specified")
	}

	dir, err := r.torrentorService.GetHLS(c.Context(), torrentInfoHash, fileHash)
	if err != nil {
		return errors.Wrap(err, "failed to get file content")
	}

	fileName := filepath.Base(c.OriginalURL())
	if fileName == "hls" {
		fileName = "output.m3u8"
	}

	return c.SendFile(filepath.Join(dir, fileName))
}

type Subtitle struct {
	Path  string
	Lang  string
	Label string
}

type Source struct {
	Path     string
	Mimetype string
}

func (r *Presentation) WatchForm(c fiber.Ctx) error {
	torrentInfoHash, err := getParamsInfoHash(c)
	if err != nil {
		return errors.Wrap(err, "failed to parse torrent")
	}

	filePath := c.Query("path")
	if filePath == "" {
		return errors.New("no file path specified")
	}

	fileMeta, err := r.torrentorService.GetFileByInfoHashAndPath(c.Context(), torrentInfoHash, filePath)
	if err != nil {
		return errors.Wrap(err, "failed to get file content")
	}

	sources := make([]Source, 0, 1)
	sources = append(sources, Source{fileMeta.Path, fileMeta.Mimetype})

	return c.Render("watch",
		fiber.Map{
			"TorrentInfoHash": torrentInfoHash.String(),
			"Path":            fileMeta.Path,
			"Mimetype":        fileMeta.Mimetype,
			"Sources":         sources,
			"Subtitles":       []Subtitle{},
		},
	)
}
