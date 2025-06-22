import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
    const [torrents, setTorrents] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetchTorrents();
    }, []);

    const fetchTorrents = async () => {
        try {
            const response = await fetch('http://localhost:8081/api/torrents');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            setTorrents(data);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    if (loading) {
        return (
            <div className="App">
                <header className="App-header">
                    <h1>Torrentor</h1>
                    <p>Loading torrents...</p>
                </header>
            </div>
        );
    }

    if (error) {
        return (
            <div className="App">
                <header className="App-header">
                    <h1>Torrentor</h1>
                    <p>Error: {error}</p>
                    <button onClick={fetchTorrents}>Retry</button>
                </header>
            </div>
        );
    }

    return (
        <div className="App">
            <header className="App-header">
                <h1>Torrentor</h1>
                <div className="torrents-container">
                    {torrents.length === 0 ? (
                        <p>No torrents found</p>
                    ) : (
                        <ul className="torrents-list">
                            {torrents.map((torrent, index) => (
                                <li key={index} className="torrent-item">
                                    {torrent.name}
                                </li>
                            ))}
                        </ul>
                    )}
                </div>
                <button onClick={fetchTorrents} className="refresh-button">
                    Refresh
                </button>
            </header>
        </div>
    );
}

export default App; 