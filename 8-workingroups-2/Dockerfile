# Stage 1: Build the Vue frontend
FROM node:24-alpine AS frontend-builder
WORKDIR /app/frontend
COPY 8-workingroups-2/frontend/package*.json ./
RUN npm install && npm install axios vis-network vis-data
COPY 8-workingroups-2/frontend/ .
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.25.8-alpine AS backend-builder

WORKDIR /app

# Copy root go.mod and common package (required for replace directive)
COPY go.mod ./
COPY common ./common

# Copy backend files
WORKDIR /app/8-workingroups-2
COPY 8-workingroups-2/go.mod 8-workingroups-2/go.sum ./
RUN go mod download
COPY 8-workingroups-2/ .
# CGO_ENABLED=0 and static flags for scratch compatibility
# -w -s flags to reduce binary size by removing debug info
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/main .

# Stage 3: Final image (Scratch - absolutely minimal)
FROM scratch
WORKDIR /
COPY --from=backend-builder /app/main /main
COPY --from=frontend-builder /app/frontend/dist /frontend/dist
EXPOSE 8080
ENTRYPOINT ["/main"]
