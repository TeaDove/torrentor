// API Configuration
export const API_CONFIG = {
    // Base URL for the API - can be overridden by environment variables
    BASE_URL: process.env.REACT_APP_API_URL || 'http://localhost:8080',
} as const;

// App Configuration
export const APP_CONFIG = {
    NAME: 'Torrentor',
    VERSION: '1.0.0',
    REFRESH_INTERVAL: 5000, // 5 seconds
} as const; 