FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.26-alpine AS backend-builder

RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/release-center .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /app/bin/release-center .
COPY --from=backend-builder /app/config.yaml.example ./config.yaml.example

RUN mkdir -p /data

EXPOSE 8080

ENV GIN_MODE=release
ENV ZENTAO_RELEASE_DATABASE_PATH=/data/release.db

ENTRYPOINT ["./release-center"]
