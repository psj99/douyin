# Tiny-DouYin
TikTok, but a homemade minimal version. **(BACKEND SERVICES ONLY!!!)**

## Requirements
- MySQL
- Redis
- 七牛OSS

## Usage
1. Copy `conf/locale/config.yaml.example` to `conf/locale/config.yaml`, then modify it as you want.

2. Migrate the: 
    ```go
    go run ./cmd/migration
    ```
3. Then
    ```go
    go run ./cmd/server
    ```

Do Not Forget go mod tidy !!!
