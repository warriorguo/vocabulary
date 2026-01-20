# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

English Wordbook - A web application for looking up word definitions and building a personal vocabulary list.

### Tech Stack
- **Frontend**: React + TypeScript + Vite
- **Backend**: Go + Gin
- **Database**: PostgreSQL
- **Dictionary API**: Free Dictionary API (api.dictionaryapi.dev)

## Project Structure

```
vocabulary/
├── frontend/                 # React app
│   ├── src/
│   │   ├── components/       # Reusable UI components
│   │   ├── pages/            # Search and Wordbook pages
│   │   ├── services/         # API client
│   │   ├── types/            # TypeScript interfaces
│   │   └── App.tsx
│   ├── package.json
│   └── vite.config.ts
├── backend/                  # Go API server
│   ├── cmd/server/main.go    # Entry point
│   ├── internal/
│   │   ├── handlers/         # HTTP handlers
│   │   ├── models/           # Data models
│   │   ├── repository/       # Database operations
│   │   └── services/         # Business logic + dict API client
│   ├── go.mod
│   └── go.sum
├── docker-compose.yml        # PostgreSQL + app services
└── CLAUDE.md
```

## Development Commands

### Backend (Go)

```bash
# Navigate to backend directory
cd backend

# Run the server (requires PostgreSQL)
go run cmd/server/main.go

# Build the binary
go build -o server ./cmd/server

# Run tests
go test ./...

# Format code
go fmt ./...
```

### Frontend (React)

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server (with API proxy to localhost:8080)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

### Docker

```bash
# Start all services (PostgreSQL, backend, frontend)
docker-compose up -d

# Start only PostgreSQL (for local development)
docker-compose up -d postgres

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/dict?word={word} | Lookup word definition (cached for 7 days) |
| GET | /api/wordbook | List saved words |
| POST | /api/wordbook | Add word to wordbook |
| DELETE | /api/wordbook/{word} | Remove word from wordbook |
| GET | /health | Health check endpoint |

## Environment Variables

### Backend
- `DATABASE_URL` - PostgreSQL connection string (default: `postgres://postgres:postgres@localhost:5432/vocabulary?sslmode=disable`)
- `PORT` - Server port (default: `8080`)

## Local Development Setup

1. Start PostgreSQL:
   ```bash
   docker-compose up -d postgres
   ```

2. Start backend (in one terminal):
   ```bash
   cd backend && go run cmd/server/main.go
   ```

3. Start frontend (in another terminal):
   ```bash
   cd frontend && npm run dev
   ```

4. Open http://localhost:5173 in your browser

## Testing

### Backend Tests

```bash
cd backend

# Run all unit tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run integration tests (requires Docker)
go test -tags integration ./...
```

### Frontend Tests

```bash
cd frontend

# Run tests in watch mode
npm run test

# Run tests once
npm run test:run

# Run tests with coverage
npm run test:coverage
```

### Manual Testing

1. Search for "hello" - verify meanings, phonetics, and audio playback
2. Click "Add to Wordbook" - verify word is added
3. Navigate to /wordbook - verify word appears with date
4. Click word to view details - verify navigation works
5. Delete word - verify removal
6. Test on mobile viewport (Chrome DevTools device toolbar)
