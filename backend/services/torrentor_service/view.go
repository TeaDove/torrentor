package torrentor_service

// func (r *Service) GetFileWithContent(
//	ctx context.Context,
//	torrentInfoHash metainfo.Hash,
//	filePath string,
// ) (schemas.FileWithContent, error) {
//	fileEnt, err := r.GetFileByInfoHashAndPath(ctx, torrentInfoHash, filePath)
//	if err != nil {
//		return schemas.FileWithContent{}, errors.Wrap(err, "failed to get torrent info")
//	}
//
//	file, err := os.Open(fileEnt.RawLocation())
//	if err != nil {
//		return schemas.FileWithContent{}, errors.Wrap(err, "error opening file")
//	}
//
//	return schemas.FileWithContent{FileEntity: fileEnt, OSFile: file}, nil
//}

// func (r *Service) GetHLS(
//	ctx context.Context,
//	torrentInfoHash metainfo.Hash,
//	fileHash string,
//	streamName string,
// ) (string, error) {
//	fileEnt, err := r.GetFileByInfoHashAndHash(ctx, torrentInfoHash, fileHash)
//	if err != nil {
//		return "", errors.Wrap(err, "failed to get torrent info")
//	}
//
//	stream, ok := fileEnt.Meta.StreamMap[streamName]
//	if !ok {
//		return "", errors.New("stream not found")
//	}
//
//	return fileEnt.LocationInUnpackAsStream(&stream, ".m3u8"), nil
//}
