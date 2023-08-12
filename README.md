# Tiny-DouYin
TikTok, but a homemade minimal version. **(BACKEND SERVICES ONLY!!!)**

## Requirements
- MySQL
- MinIO (or other object storage)

## Usage
1. Copy `conf/locale/config.yaml.example` to `conf/locale/config.yaml`, then modify it as you want.
2. Migrate the database: 
    ```go
    go run ./cmd/migration
    ```
3. Then
    ```go
    go run ./cmd/server
    ```

