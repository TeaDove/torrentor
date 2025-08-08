import { API_CONFIG } from '../config';
import {
    Torrent,
    TorrentEntity,
    ApiResponse,
    AddTorrentRequest,
    DownloadTorrentRequest,
    TorrentActionResponse,
    Stats,
    API_ENDPOINTS
} from './models';

// Helper function to build full API URLs
const buildApiUrl = (endpoint: string): string => {
    return `${API_CONFIG.BASE_URL}${endpoint}`;
};

// Request configuration
const requestConfig = {
    headers: { 'Content-Type': 'application/json' },
    timeout: 10000, // 10 seconds
};

// Generic API request function
const apiRequest = async <T>(
    endpoint: string,
    options: RequestInit = {}
): Promise<T> => {
    const url = buildApiUrl(endpoint);
    const config = {
        ...requestConfig,
        ...options,
        headers: {
            ...requestConfig.headers,
            ...options.headers,
        },
    };

    const response = await fetch(url, config);

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
};

// Torrent API functions
export const torrentApi = {
    // Get all torrents
    getAll: async (): Promise<Torrent[]> => {
        return apiRequest<Torrent[]>(API_ENDPOINTS.TORRENTS);
    },

    // Get torrent by ID
    getById: async (id: string): Promise<Torrent> => {
        return apiRequest<Torrent>(API_ENDPOINTS.TORRENT_DETAILS(id));
    },

    // Get torrent details with files by infoHash
    getTorrentDetails: async (infoHash: string): Promise<TorrentEntity> => {
        return apiRequest<TorrentEntity>(API_ENDPOINTS.TORRENT_DETAILS(infoHash));
    },

    // Add new torrent
    add: async (request: AddTorrentRequest): Promise<TorrentActionResponse> => {
        const formData = new FormData();

        if (request.magnet) {
            formData.append('magnet', request.magnet);
        }

        if (request.file) {
            formData.append('file', request.file);
        }

        if (request.path) {
            formData.append('path', request.path);
        }

        const url = buildApiUrl(API_ENDPOINTS.ADD_TORRENT);
        const response = await fetch(url, {
            method: 'POST',
            body: formData,
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return response.json();
    },

    // Download torrent with magnet link
    download: async (request: DownloadTorrentRequest): Promise<Torrent> => {
        return apiRequest<Torrent>(
            API_ENDPOINTS.DOWNLOAD_TORRENT,
            {
                method: 'POST',
                body: JSON.stringify(request),
            }
        );
    },

    // Delete torrent
    delete: async (infoHash: string): Promise<TorrentActionResponse> => {
        return apiRequest<TorrentActionResponse>(
            API_ENDPOINTS.DELETE_TORRENT(infoHash),
            { method: 'DELETE' }
        );
    },

    // Pause torrent
    pause: async (id: string): Promise<TorrentActionResponse> => {
        return apiRequest<TorrentActionResponse>(
            API_ENDPOINTS.PAUSE_TORRENT(id),
            { method: 'POST' }
        );
    },

    // Resume torrent
    resume: async (id: string): Promise<TorrentActionResponse> => {
        return apiRequest<TorrentActionResponse>(
            API_ENDPOINTS.RESUME_TORRENT(id),
            { method: 'POST' }
        );
    },

    // Get stats
    getStats: async (): Promise<Stats> => {
        return apiRequest<Stats>(API_ENDPOINTS.STATS);
    },

    // Download file
    downloadFile: async (infoHash: string, pathHash: string): Promise<Blob> => {
        const url = buildApiUrl(API_ENDPOINTS.DOWNLOAD_FILE(infoHash, pathHash));
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                ...requestConfig.headers,
            },
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return response.blob();
    },
};

// Export the main API object
export default torrentApi; 