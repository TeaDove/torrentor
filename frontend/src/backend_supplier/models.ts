// Backend API Models

export interface Meta {
    pieces: number;
    piecesLength: number;
    magnet: string;
}

export interface FileEntity {
    name: string;
    path: string;
    pathHash: string;
    mimetype?: string;
    size: number;
    completed: boolean;
    meta?: any; // ffmpeg_service.Metadata
}

export interface TorrentEntity {
    createdAt: string; // ISO date string from Go's time.Time
    name: string;
    filePathMap?: { [key: string]: FileEntity };
    size: number; // Byte size from conv_utils.Byte
    infoHash: string; // metainfo.Hash as string
    completed: boolean;
    meta?: Meta;
}

// Legacy Torrent interface for backward compatibility
export interface Torrent {
    createdAt: string; // ISO date string from Go's time.Time
    name: string;
    size: number; // Byte size from conv_utils.Byte
    infoHash: string; // metainfo.Hash as string
    completed: boolean;
    meta?: Meta; // Optional Meta struct
}

export interface ApiResponse<T> {
    data: T;
    error?: string;
    message?: string;
}

export interface AddTorrentRequest {
    magnet?: string;
    file?: File;
    path?: string;
}

export interface DownloadTorrentRequest {
    magnet: string;
}

export interface TorrentActionResponse {
    success: boolean;
    message?: string;
    error?: string;
}

export interface Stats {
    [key: string]: any; // Flexible stats object
}

// API Endpoints
export const API_ENDPOINTS = {
    TORRENTS: '/api/torrents',
    TORRENT_DETAILS: (id: string) => `/api/torrents/${id}`,
    ADD_TORRENT: '/api/torrents/add',
    DOWNLOAD_TORRENT: '/api/torrents/download',
    DELETE_TORRENT: (infoHash: string) => `/api/torrents/${infoHash}`,
    PAUSE_TORRENT: (id: string) => `/api/torrents/${id}/pause`,
    RESUME_TORRENT: (id: string) => `/api/torrents/${id}/resume`,
    STATS: '/api/stats',
    DOWNLOAD_FILE: (infoHash: string, pathHash: string) => `/api/torrents/${infoHash}/files/${pathHash}`,
} as const; 