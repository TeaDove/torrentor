import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { APP_CONFIG } from './config';
import { torrentApi, Torrent, Stats } from './backend_supplier';

const Home: React.FC = () => {
    const navigate = useNavigate();
    const [torrents, setTorrents] = useState<Torrent[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [stats, setStats] = useState<Stats | null>(null);
    const [statsError, setStatsError] = useState<string | null>(null);
    const [magnetLink, setMagnetLink] = useState<string>('');
    const [downloadLoading, setDownloadLoading] = useState<boolean>(false);
    const [downloadMessage, setDownloadMessage] = useState<string | null>(null);
    const [deletingTorrents, setDeletingTorrents] = useState<Set<string>>(new Set());

    useEffect(() => {
        fetchTorrents();
        fetchStats();
    }, []);

    // Fetch stats every 10 seconds
    useEffect(() => {
        const statsInterval = setInterval(() => {
            fetchStats();
        }, 10000);

        return () => clearInterval(statsInterval);
    }, []);

    const fetchTorrents = async (): Promise<void> => {
        try {
            const data = await torrentApi.getAll();
            setTorrents(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An unknown error occurred');
        } finally {
            setLoading(false);
        }
    };

    const fetchStats = async (): Promise<void> => {
        try {
            const data = await torrentApi.getStats();
            setStats(data);
            setStatsError(null);
        } catch (err) {
            setStatsError(err instanceof Error ? err.message : 'Failed to fetch stats');
        }
    };

    const handleDownload = async (): Promise<void> => {
        if (!magnetLink.trim()) {
            setDownloadMessage('Please enter a magnet link');
            return;
        }

        setDownloadLoading(true);
        setDownloadMessage(null);

        try {
            const torrent = await torrentApi.download({ magnet: magnetLink.trim() });
            setDownloadMessage(`Torrent "${torrent.name}" added successfully!`);
            setMagnetLink('');
            // Refresh torrents list after successful download
            setTimeout(() => {
                fetchTorrents();
            }, 1000);
        } catch (err) {
            setDownloadMessage(err instanceof Error ? err.message : 'Failed to add torrent');
        } finally {
            setDownloadLoading(false);
        }
    };

    const handleDelete = async (torrent: Torrent): Promise<void> => {
        setDeletingTorrents(prev => new Set(prev).add(torrent.infoHash));

        try {
            await torrentApi.delete(torrent.infoHash);
            // Remove the torrent from the local state immediately
            setTorrents(prev => prev.filter(t => t.infoHash !== torrent.infoHash));
        } catch (err) {
            console.error('Failed to delete torrent:', err);
            // Optionally show an error message to the user
        } finally {
            setDeletingTorrents(prev => {
                const newSet = new Set(prev);
                newSet.delete(torrent.infoHash);
                return newSet;
            });
        }
    };

    const handleKeyPress = (e: React.KeyboardEvent): void => {
        if (e.key === 'Enter') {
            handleDownload();
        }
    };

    const handleTorrentClick = (torrent: Torrent): void => {
        navigate(`/torrents/${torrent.infoHash}`);
    };

    const handleDeleteClick = (e: React.MouseEvent, torrent: Torrent): void => {
        e.stopPropagation(); // Prevent navigation when clicking delete button
        handleDelete(torrent);
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
                    <h1>{APP_CONFIG.NAME}</h1>
                    <p>Loading torrents...</p>
                </header>
            </div>
        );
    }

    if (error) {
        return (
            <div className="App">
                <header className="App-header">
                    <h1>{APP_CONFIG.NAME}</h1>
                    <p>Error: {error}</p>
                    <button onClick={fetchTorrents}>Retry</button>
                </header>
            </div>
        );
    }

    return (
        <div className="App">
            <header className="App-header">
                <h1>{APP_CONFIG.NAME}</h1>

                {/* Magnet Link Input Section */}
                <div className="magnet-input-section">
                    <div className="input-group">
                        <input
                            type="text"
                            placeholder="Enter magnet link..."
                            value={magnetLink}
                            onChange={(e) => setMagnetLink(e.target.value)}
                            onKeyPress={handleKeyPress}
                            className="magnet-input"
                            disabled={downloadLoading}
                        />
                        <button
                            onClick={handleDownload}
                            disabled={downloadLoading || !magnetLink.trim()}
                            className="download-button"
                        >
                            {downloadLoading ? 'Adding...' : 'Add'}
                        </button>
                    </div>
                    {downloadMessage && (
                        <div className={`download-message ${downloadMessage.includes('successfully') ? 'success' : 'error'}`}>
                            {downloadMessage}
                        </div>
                    )}
                </div>

                <div className="torrents-container">
                    {torrents.length === 0 ? (
                        <p>No torrents found</p>
                    ) : (
                        <ul className="torrents-list">
                            {torrents.map((torrent: Torrent, index: number) => (
                                <li
                                    key={torrent.infoHash || index}
                                    className="torrent-item clickable"
                                    onClick={() => handleTorrentClick(torrent)}
                                >
                                    <div className="torrent-content">
                                        <div className="torrent-name">
                                            {torrent.name}
                                        </div>
                                        <div className="torrent-info">
                                            <div className="torrent-size">
                                                Size: {formatFileSize(torrent.size)}
                                            </div>
                                            <div className="torrent-status">
                                                Status: {torrent.completed ? 'Completed' : 'In Progress'}
                                            </div>
                                            <div className="torrent-created">
                                                Added: {formatDate(torrent.createdAt)}
                                            </div>
                                            {torrent.meta && (
                                                <div className="torrent-meta">
                                                    <div>Pieces: {torrent.meta.pieces}</div>
                                                    <div>Piece Length: {formatFileSize(torrent.meta.piecesLength)}</div>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                    <div className="torrent-actions">
                                        <button
                                            onClick={(e) => handleDeleteClick(e, torrent)}
                                            disabled={deletingTorrents.has(torrent.infoHash)}
                                            className="delete-button"
                                            title="Delete torrent"
                                        >
                                            {deletingTorrents.has(torrent.infoHash) ? 'Deleting...' : 'Delete'}
                                        </button>
                                    </div>
                                </li>
                            ))}
                        </ul>
                    )}
                </div>
                <button onClick={fetchTorrents} className="refresh-button">
                    Refresh
                </button>
            </header>

            {/* Stats display in bottom left corner */}
            <div className="stats-container">
                {statsError ? (
                    <div className="stats-error">
                        Stats Error: {statsError}
                    </div>
                ) : stats ? (
                    <div className="stats-content">
                        <h3>Services Stats</h3>
                        <pre>{JSON.stringify(stats, null, 2)}</pre>
                    </div>
                ) : (
                    <div className="stats-loading">
                        Loading stats...
                    </div>
                )}
            </div>
        </div>
    );
};

export default Home; 