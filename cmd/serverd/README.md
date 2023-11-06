# About

serverd is HTTP Server interface implementation of wallet-scanner service.

# Usage

Please run program with `--help` flag to see available configuration options if running manually.

# Build and run with docker

Copy `.env.example` to `.env` to the project root and update configuration values:

* Set `RISKPROVIDER_BLOCKMATE_APIKEY` value to a valid token

Run application:

```bash
docker compose up
```

# Test

Screen ethereum wallet address:

```bash
curl -X POST 'http://localhost/wallet/0xe9e9afac38e64728f1afbb2b65dec7be7c704c05/categories' -v
```

Retrieve history of risk categories for the ethereum address:

```bash
curl 'http://localhost/wallet/0xe9e9afac38e64728f1afbb2b65dec7be7c704c05/categories' -v
```
