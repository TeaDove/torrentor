# Torrentor Frontend

This is the TypeScript React frontend for the Torrentor application.

## Features

- **TypeScript**: Full TypeScript support with strict type checking
- **React 18**: Latest React features with hooks
- **Modern UI**: Clean, responsive design with gradient backgrounds
- **Real-time Updates**: Displays torrent information with progress, status, and size
- **Error Handling**: Robust error handling with retry functionality
- **Modular API**: Centralized backend supplier for all API requests
- **Live Stats**: Real-time service statistics displayed in bottom left corner
- **Magnet Link Download**: Add torrents via magnet links with real-time feedback
- **Torrent Management**: Delete torrents with individual delete buttons

## Project Structure

```
src/
├── App.tsx                    # Main application component
├── index.tsx                  # Application entry point
├── config.ts                  # Configuration (base URL only)
├── App.css                    # Application styles
├── react-app-env.d.ts        # TypeScript declarations
└── backend_supplier/         # Backend API integration
    ├── index.ts              # Main exports
    ├── models.ts             # Backend data models
    └── api.ts               # API client and requests
```

## Configuration

### API Configuration

The application uses a simplified configuration system with only the base URL:

1. **Environment Variable** (Recommended):
   ```bash
   # Create a .env file in the frontend directory
   REACT_APP_API_URL=http://localhost:8080
   ```

2. **Default Configuration**:
   - Default API URL: `http://localhost:8080`
   - Can be overridden using environment variables

## Backend Supplier

The `backend_supplier` folder contains all backend-related code:

### Models (`models.ts`)
- `Torrent`: Interface matching Go backend structure
- `Meta`: Optional metadata for torrents
- `ApiResponse<T>`: Generic API response wrapper
- `AddTorrentRequest`: Request interface for adding torrents
- `DownloadTorrentRequest`: Request interface for downloading torrents
- `TorrentActionResponse`: Response interface for torrent actions
- `Stats`: Interface for service statistics
- `API_ENDPOINTS`: All API endpoint definitions

### Torrent Model Structure
The frontend `Torrent` interface matches the Go backend `TorrentEntity`:

```typescript
interface Torrent {
    createdAt: string;     // ISO date string from Go's time.Time
    name: string;          // Torrent name
    size: number;          // Byte size from conv_utils.Byte
    infoHash: string;      // metainfo.Hash as string
    completed: boolean;    // Download completion status
    meta?: Meta;          // Optional metadata
}

interface Meta {
    pieces: number;        // Number of pieces
    piecesLength: number;  // Length of each piece
    magnet: string;        // Magnet link
}
```

### API Client (`api.ts`)
- `torrentApi`: Main API client with all torrent operations
- Generic request handling with error management
- Type-safe API calls

### Available API Methods
- `torrentApi.getAll()` - Get all torrents (returns `Torrent[]`)
- `torrentApi.getById(id)` - Get torrent details (returns `Torrent`)
- `torrentApi.add(request)` - Add new torrent (returns `TorrentActionResponse`)
- `torrentApi.download(request)` - Download torrent with magnet link (returns `Torrent`)
- `torrentApi.delete(infoHash)` - Delete torrent (returns `TorrentActionResponse`)
- `torrentApi.pause(id)` - Pause torrent (returns `TorrentActionResponse`)
- `torrentApi.resume(id)` - Resume torrent (returns `TorrentActionResponse`)
- `torrentApi.getStats()` - Get service statistics (returns `Stats`)

## Features

### Torrent Management
Each torrent item displays:
- **Torrent Name**: Primary identifier
- **File Size**: Formatted size display (Bytes, KB, MB, GB, TB)
- **Completion Status**: "Completed" or "In Progress"
- **Creation Date**: When the torrent was added
- **Metadata**: Optional piece information and magnet link
- **Delete Button**: Individual delete button for each torrent
- **Loading States**: Button shows "Deleting..." during request
- **Immediate UI Update**: Torrent removed from list immediately after deletion
- **Error Handling**: Graceful error handling for failed deletions

### Magnet Link Download
The application includes a magnet link input section:
- **Input Field**: Paste magnet links directly
- **Add Button**: Triggers download with visual feedback
- **Enter Key**: Press Enter to submit
- **Real-time Feedback**: Success/error messages with torrent name
- **Auto-refresh**: Torrent list updates after successful download
- **Loading States**: Button shows "Adding..." during request

### Live Statistics
The application displays real-time service statistics in the bottom left corner:
- **Auto-refresh**: Updates every 10 seconds
- **JSON Display**: Shows raw JSON response from `/api/stats`
- **Error Handling**: Displays error messages if stats fetch fails
- **Styling**: Dark overlay with monospace font for readability

## Development

### Prerequisites

- Node.js (v14 or higher)
- npm

### Installation

```bash
npm install
```

### Environment Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Modify the `.env` file to set your API URL:
   ```bash
   REACT_APP_API_URL=http://localhost:8080
   ```

### Running the Development Server

```bash
npm start
```

The application will be available at `http://localhost:3000`.

### Building for Production

```bash
npm run build
```

### Running Tests

```bash
npm test
```

## API Integration

The frontend uses the `backend_supplier` module to interact with the backend API. The default endpoint is `http://localhost:8080/api/torrents`.

### Configuration Options

- **Development**: `http://localhost:8080`
- **Production**: Set via `REACT_APP_API_URL` environment variable
- **Custom**: Modify `src/config.ts` for base URL changes

## Styling

The application uses CSS modules with a modern gradient design and responsive layout. The UI includes:

- Gradient background
- Glassmorphism effects
- Hover animations
- Progress indicators
- Status displays
- Live stats overlay
- Magnet link input section
- Delete buttons with hover effects
- Metadata display with styled containers

## TypeScript Benefits

- **Type Safety**: Catch errors at compile time
- **Better IDE Support**: Enhanced autocomplete and refactoring
- **Self-Documenting Code**: Types serve as documentation
- **Easier Maintenance**: Clear interfaces and contracts
- **Modular Architecture**: Separated concerns with backend supplier 