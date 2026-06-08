# Stage 1: Build the Vue frontend
FROM node:24-alpine AS frontend-builder
WORKDIR /app/frontend
COPY 7-workingroups-1/frontend/package*.json ./
RUN npm install && npm install axios
COPY 7-workingroups-1/frontend/ .
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.25.8-alpine AS backend-builder

WORKDIR /app

# Copy root go.mod and common package (required for replace directive)
COPY go.mod ./
COPY common ./common

# Copy backend files
WORKDIR /app/7-workingroups-1
COPY 7-workingroups-1/go.mod 7-workingroups-1/go.sum ./
RUN go mod download
COPY 7-workingroups-1/ .
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
