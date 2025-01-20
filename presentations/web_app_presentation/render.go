package web_app_presentation

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
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

func getParamsFileHash(c fiber.Ctx) (string, error) {
	hash := c.Params("filehash")
	if hash == "" {
		return "", errors.New("no filehash")
	}

	return hash, nil
}

func getParamsStream(c fiber.Ctx) (string, error) {
	streamName := c.Params("name")
	if streamName == "" {
		return "", errors.New("no stream name specified")
	}

	return streamName, nil
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

	return c.SendFile(file.Location())
}

func (r *Presentation) HLSForm(c fiber.Ctx) error {
	torrentInfoHash, err := getParamsInfoHash(c)
	if err != nil {
		return errors.Wrap(err, "failed to parse torrent")
	}

	fileHash, err := getParamsFileHash(c)
	if err != nil {
		return errors.New("bad file hash")
	}

	file, err := r.torrentorService.GetFileByInfoHashAndHash(c.Context(), torrentInfoHash, fileHash)
	if err != nil {
		return errors.Wrap(err, "failed to get file content")
	}

	fileName := filepath.Base(c.OriginalURL())
	if strings.HasSuffix(fileName, ".ts") {
		fileName = "output.m3u8"
	}

	return c.SendFile(filepath.Join(file.LocationInUnpack(), fileName))
}

type Subtitle struct {
	Path  string
	Lang  string
	Label string
}

type Source struct {
	StreamName string
}

func (r *Presentation) WatchForm(c fiber.Ctx) error {
	torrentInfoHash, err := getParamsInfoHash(c)
	if err != nil {
		return errors.Wrap(err, "failed to parse torrent info hash")
	}

	fileHash, err := getParamsFileHash(c)
	if err != nil {
		return errors.New("bad file hash")
	}

	fileMeta, err := r.torrentorService.GetFileByInfoHashAndHash(c.Context(), torrentInfoHash, fileHash)
	if err != nil {
		return errors.Wrap(err, "failed to get file content")
	}

	streamName, err := getParamsStream(c)
	if err != nil {
		return errors.Wrap(err, "no stream name specified")
	}

	stream, ok := fileMeta.Meta.StreamMap[streamName]
	if !ok {
		return errors.New("no such stream")
	}

	sources := []Source{{StreamName: stream.StreamFile(".m3u8")}}

	return c.Render("watch",
		fiber.Map{
			"TorrentInfoHash": torrentInfoHash.String(),
			"Path":            fileMeta.Path,
			"Mimetype":        fileMeta.Mimetype,
			"FileHash":        fileHash,
			"Sources":         sources,
			"Subtitles":       []Subtitle{},
		},
	)
}
