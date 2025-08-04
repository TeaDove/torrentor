import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';
import Home from './Home';
import TorrentDetails from './TorrentDetails';

const App: React.FC = () => {
    return (
        <Router>
            <div className="App">
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/torrents/:infoHash" element={<TorrentDetails />} />
                </Routes>
            </div>
        </Router>
    );
};

export default App; 