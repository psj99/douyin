# Tiny-DouYin
TikTok, but a homemade minimal version. **(BACKEND SERVICES ONLY!!!)**

## Requirements
- MySQL
- MinIO (or other object storage)
- FFmpeg
- QiNiu SDK

## Usage
1. Copy `conf/locale/config.yaml.example` to `conf/locale/config.yaml`, then modify it as you want.
2. Run `go run ./cmd/migration` to initialize database
3. Use `go run ./cmd/server` to start the server
