package web_app_presentation

import (
	"path/filepath"
	"torrentor/services/ffmpeg_service"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
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
	//torrentInfoHash, err := getParamsInfoHash(c)
	//if err != nil {
	//	return errors.Wrap(err, "failed to parse torrent")
	//}
	//
	//fileHash, err := getParamsFileHash(c)
	//if err != nil {
	//	return errors.New("bad file hash")
	//}
	//
	//file, err := r.torrentorService.GetFileByInfoHashAndHash(c.Context(), torrentInfoHash, fileHash)
	//if err != nil {
	//	return errors.Wrap(err, "failed to get file content")
	//}

	return c.SendFile(filepath.Join("./data", c.OriginalURL()), fiber.SendFile{ByteRange: true})
}

type Subtitle struct {
	StreamName string
	Lang       string
	Label      string
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

	err = r.torrentorService.UnpackIfNeeded(c.Context(), fileMeta)
	if err != nil {
		return errors.Wrap(err, "failed to unpack file")
	}

	streamName, err := getParamsStream(c)
	if err != nil {
		return errors.Wrap(err, "no stream name specified")
	}

	stream, ok := fileMeta.Meta.StreamMap[streamName]
	if !ok {
		return errors.New("no such stream")
	}

	sources := []Source{{StreamName: stream.String()}}

	subtitles := make([]Subtitle, 0)
	for streamName, stream = range fileMeta.Meta.StreamMap {
		if stream.CodecType == ffmpeg_service.CodecTypeSubtitle {
			subtitles = append(subtitles, Subtitle{
				StreamName: streamName,
				Lang:       stream.Tags.Language,
				Label:      stream.Tags.Title,
			})
		}
	}

	return c.Render("watch",
		fiber.Map{
			"TorrentInfoHash": torrentInfoHash.String(),
			"Path":            fileMeta.Path,
			"Mimetype":        fileMeta.Mimetype,
			"FileHash":        fileHash,
			"Sources":         sources,
			"Subtitles":       subtitles,
		},
	)
}

//
//func (r *Presentation) HLSForm(c fiber.Ctx) error {
//	torrentInfoHash, err := getParamsInfoHash(c)
//	if err != nil {
//		return errors.Wrap(err, "failed to parse torrent")
//	}
//
//	fileHash, err := getParamsFileHash(c)
//	if err != nil {
//		return errors.New("bad file hash")
//	}
//
//	file, err := r.torrentorService.GetFileByInfoHashAndHash(c.Context(), torrentInfoHash, fileHash)
//	if err != nil {
//		return errors.Wrap(err, "failed to get file content")
//	}
//
//	streamName, err := getParamsStream(c)
//	if err != nil {
//		return errors.Wrap(err, "no stream name specified")
//	}
//
//	stream, ok := file.Meta.StreamMap[streamName]
//	if !ok {
//		return errors.New("no such stream")
//	}
//
//	fileName := fmt.Sprintf(".m3u8/%s", filepath.Base(c.OriginalURL()))
//
//	return c.SendFile(file.LocationInUnpackAsStream(&stream, fileName))
//}
