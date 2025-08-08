import React, { useState, useEffect, useRef } from 'react';
import { useParams, Link } from 'react-router-dom';
import { torrentApi, TorrentEntity, FileEntity } from './backend_supplier';
import { API_CONFIG } from './config';
import './App.css';

interface TreeNode {
    name: string;
    path: string;
    pathHash?: string;
    size: number;
    completed: boolean;
    mimetype?: string;
    children: { [key: string]: TreeNode };
    isDirectory: boolean;
}

interface VideoPlayerProps {
    infoHash: string;
    pathHash: string;
    fileName: string;
    onClose: () => void;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ infoHash, pathHash, fileName, onClose }) => {
    const videoRef = useRef<HTMLVideoElement>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const videoUrl = `${API_CONFIG.BASE_URL}/api/torrents/${infoHash}/files/${pathHash}`;

    const handleVideoLoad = () => {
        setIsLoading(false);
    };

    const handleVideoError = () => {
        setError('Failed to load video');
        setIsLoading(false);
    };

    return (
        <div className="video-player-overlay" onClick={onClose}>
            <div className="video-player-container" onClick={(e) => e.stopPropagation()}>
                <div className="video-player-header">
                    <h3>{fileName}</h3>
                    <button className="video-player-close" onClick={onClose}>
                        ✕
                    </button>
                </div>
                <div className="video-player-content">
                    {isLoading && (
                        <div className="video-loading">
                            <div className="loading-spinner"></div>
                            <p>Loading video...</p>
                        </div>
                    )}
                    {error && (
                        <div className="video-error">
                            <p>{error}</p>
                        </div>
                    )}
                    <video
                        ref={videoRef}
                        controls
                        preload="metadata"
                        onLoadedData={handleVideoLoad}
                        onError={handleVideoError}
                        style={{ display: isLoading || error ? 'none' : 'block' }}
                    >
                        <source src={videoUrl} type="video/webm" />
                        <source src={videoUrl} type="video/mp4" />
                        <source src={videoUrl} type="video/ogg" />
                        Your browser does not support the video tag.
                    </video>
                </div>
            </div>
        </div>
    );
};

const TorrentDetails: React.FC = () => {
    const { infoHash } = useParams<{ infoHash: string }>();
    const [torrent, setTorrent] = useState<TorrentEntity | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [videoPlayer, setVideoPlayer] = useState<{
        infoHash: string;
        pathHash: string;
        fileName: string;
    } | null>(null);

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

    // Helper function to build tree structure from files
    const buildFileTree = (files: FileEntity[]): TreeNode[] => {
        const root: { [key: string]: TreeNode } = {};

        files.forEach(file => {
            const pathParts = file.path.split('/').filter(part => part.length > 0);
            let currentLevel = root;

            // Handle files in root directory
            if (pathParts.length === 1) {
                // File is in root directory
                const fileName = pathParts[0];
                currentLevel[fileName] = {
                    name: fileName,
                    path: file.path,
                    pathHash: file.pathHash,
                    size: file.size,
                    completed: file.completed,
                    mimetype: file.mimetype,
                    children: {},
                    isDirectory: false
                };
            } else {
                // Create directory nodes for all path parts except the last one
                for (let i = 0; i < pathParts.length - 1; i++) {
                    const part = pathParts[i];
                    if (!currentLevel[part]) {
                        currentLevel[part] = {
                            name: part,
                            path: pathParts.slice(0, i + 1).join('/'),
                            size: 0,
                            completed: true,
                            children: {},
                            isDirectory: true
                        };
                    }
                    currentLevel = currentLevel[part].children;
                }

                // Add the file node
                const fileName = pathParts[pathParts.length - 1];
                currentLevel[fileName] = {
                    name: fileName,
                    path: file.path,
                    pathHash: file.pathHash,
                    size: file.size,
                    completed: file.completed,
                    mimetype: file.mimetype,
                    children: {},
                    isDirectory: false
                };
            }
        });

        return Object.values(root);
    };

    // Helper function to check if a file is a video
    const isVideoFile = (mimetype?: string): boolean => {
        if (!mimetype) return false;
        return mimetype.startsWith('video/');
    };

    // Helper function to handle video playback
    const handleWatchVideo = (pathHash: string, fileName: string) => {
        if (!infoHash) return;
        setVideoPlayer({
            infoHash,
            pathHash,
            fileName
        });
    };

    // Helper function to close video player
    const closeVideoPlayer = () => {
        setVideoPlayer(null);
    };

    // Helper function to handle file download
    const handleDownload = async (pathHash: string, fileName: string) => {
        if (!infoHash || !pathHash) return;

        try {
            const blob = await torrentApi.downloadFile(infoHash, pathHash);

            // Create a download link
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = fileName;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            window.URL.revokeObjectURL(url);
        } catch (err) {
            console.error('Download failed:', err);
            alert('Download failed. Please try again.');
        }
    };

    // Helper function to render tree nodes recursively
    const renderTreeNode = (node: TreeNode, level: number = 0): JSX.Element => {
        const indent = level * 20; // 20px per level
        const isDirectory = node.isDirectory;

        return (
            <div key={node.path} className="tree-node">
                <div
                    className={`tree-item ${isDirectory ? 'directory' : 'file'}`}
                    style={{ paddingLeft: `${indent}px` }}
                >
                    <div className="tree-item-content">
                        <span className="tree-item-name">
                            {isDirectory ? '📁 ' : '📄 '}{node.name}
                        </span>
                        {!isDirectory && (
                            <div className="tree-item-info">
                                <span className="tree-item-size">{formatFileSize(node.size)}</span>
                                {!node.completed && (
                                    <span className="tree-item-status in-progress">In Progress</span>
                                )}
                                {node.mimetype && (
                                    <span className="tree-item-mimetype">{node.mimetype}</span>
                                )}
                                {node.pathHash && (
                                    <button
                                        className="download-file-button"
                                        onClick={() => handleDownload(node.pathHash!, node.name)}
                                        disabled={!node.completed}
                                        title={node.completed ? 'Download file' : 'File not completed'}
                                    >
                                        ⬇️ Download
                                    </button>
                                )}
                                {isVideoFile(node.mimetype) && (
                                    <button
                                        className="watch-file-button"
                                        onClick={() => handleWatchVideo(node.pathHash!, node.name)}
                                        disabled={!node.completed}
                                        title={node.completed ? 'Watch file' : 'File not completed'}
                                    >
                                        ▶️ Watch
                                    </button>
                                )}
                            </div>
                        )}
                    </div>
                </div>
                {isDirectory && Object.values(node.children).length > 0 && (
                    <div className="tree-children">
                        {Object.values(node.children)
                            .sort((a, b) => {
                                // Sort directories first, then files
                                if (a.isDirectory && !b.isDirectory) return -1;
                                if (!a.isDirectory && b.isDirectory) return 1;
                                return a.name.localeCompare(b.name);
                            })
                            .map(child => renderTreeNode(child, level + 1))}
                    </div>
                )}
            </div>
        );
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
    const fileTree = buildFileTree(files);

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
                        <div className="files-tree">
                            {fileTree
                                .sort((a, b) => {
                                    // Sort directories first, then files
                                    if (a.isDirectory && !b.isDirectory) return -1;
                                    if (!a.isDirectory && b.isDirectory) return 1;
                                    return a.name.localeCompare(b.name);
                                })
                                .map(node => renderTreeNode(node))}
                        </div>
                    )}
                </div>
            </header>
            {videoPlayer && (
                <VideoPlayer
                    infoHash={videoPlayer.infoHash}
                    pathHash={videoPlayer.pathHash}
                    fileName={videoPlayer.fileName}
                    onClose={closeVideoPlayer}
                />
            )}
        </div>
    );
};

export default TorrentDetails; 