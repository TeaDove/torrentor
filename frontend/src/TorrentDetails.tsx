import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { torrentApi, TorrentEntity, FileEntity } from './backend_supplier';
import './App.css';

const TorrentDetails: React.FC = () => {
    const { infoHash } = useParams<{ infoHash: string }>();
    const [torrent, setTorrent] = useState<TorrentEntity | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (infoHash) {
            fetchTorrentDetails();
        }
    }, [infoHash]);

    const fetchTorrentDetails = async (): Promise<void> => {
        if (!infoHash) return;

        try {
            setLoading(true);
            const data = await torrentApi.getTorrentDetails(infoHash);
            setTorrent(data);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to fetch torrent details');
        } finally {
            setLoading(false);
        }
    };

    // Helper function to format file size
    const formatFileSize = (bytes: number): string => {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    // Helper function to format date
    const formatDate = (dateString: string): string => {
        return new Date(dateString).toLocaleString();
    };

    if (loading) {
        return (
            <div className="App">
                <header className="App-header">
                    <h1>Loading torrent details...</h1>
                </header>
            </div>
        );
    }

    if (error) {
        return (
            <div className="App">
                <header className="App-header">
                    <h1>Error</h1>
                    <p>{error}</p>
                    <Link to="/" className="back-button">Back to Home</Link>
                </header>
            </div>
        );
    }

    if (!torrent) {
        return (
            <div className="App">
                <header className="App-header">
                    <h1>Torrent not found</h1>
                    <Link to="/" className="back-button">Back to Home</Link>
                </header>
            </div>
        );
    }

    const files = torrent.filePathMap ? Object.values(torrent.filePathMap) : [];

    return (
        <div className="App">
            <header className="App-header">
                <div className="torrent-details-header">
                    <Link to="/" className="back-button">← Back to Home</Link>
                    <h1>{torrent.name}</h1>
                </div>

                <div className="torrent-info-section">
                    <div className="torrent-info-grid">
                        <div className="info-item">
                            <strong>Status:</strong> {torrent.completed ? 'Completed' : 'In Progress'}
                        </div>
                        <div className="info-item">
                            <strong>Size:</strong> {formatFileSize(torrent.size)}
                        </div>
                        <div className="info-item">
                            <strong>Added:</strong> {formatDate(torrent.createdAt)}
                        </div>
                        <div className="info-item">
                            <strong>Info Hash:</strong> {torrent.infoHash}
                        </div>
                        {torrent.meta && (
                            <>
                                <div className="info-item">
                                    <strong>Pieces:</strong> {torrent.meta.pieces}
                                </div>
                                <div className="info-item">
                                    <strong>Piece Length:</strong> {formatFileSize(torrent.meta.piecesLength)}
                                </div>
                                <div className="info-item">
                                    <strong>Magnet:</strong>
                                    <a href={torrent.meta.magnet} className="magnet-link" target="_blank" rel="noopener noreferrer">
                                        Copy Link
                                    </a>
                                </div>
                            </>
                        )}
                    </div>
                </div>

                <div className="files-section">
                    <h2>Files ({files.length})</h2>
                    {files.length === 0 ? (
                        <p>No files found</p>
                    ) : (
                        <div className="files-list">
                            {files.map((file: FileEntity, index: number) => (
                                <div key={file.pathHash || index} className="file-item">
                                    <div className="file-content">
                                        <div className="file-name">{file.name}</div>
                                        <div className="file-path">{file.path}</div>
                                        <div className="file-info">
                                            <span className="file-size">{formatFileSize(file.size)}</span>
                                            <span className={`file-status ${file.completed ? 'completed' : 'in-progress'}`}>
                                                {file.completed ? 'Completed' : 'In Progress'}
                                            </span>
                                            {file.mimetype && (
                                                <span className="file-mimetype">{file.mimetype}</span>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </header>
        </div>
    );
};

export default TorrentDetails; 