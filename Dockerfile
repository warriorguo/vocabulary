# ============================================
# Stage 1: Build Frontend
# ============================================
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./
RUN npm ci

# Copy source and build
COPY frontend/ ./
RUN npm run build

# ============================================
# Stage 2: Build Backend
# ============================================
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app/backend

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source and build
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# ============================================
# Stage 3: Runtime
# ============================================
FROM alpine:3.19

RUN apk --no-cache add ca-certificates nginx supervisor

WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /server ./server

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/http.d/default.conf

# Copy supervisor config
COPY supervisord.conf /etc/supervisord.conf

# Expose port
EXPOSE 80

# Start supervisor (manages both nginx and backend)
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
