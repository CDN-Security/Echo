# HTTP Echo Server

## Installation

```bash
go install github.com/CDN-Security/Echo/cmd/http-echo-server@latest
```

## Usage

```bash
# Download config file
wget https://github.com/CDN-Security/Echo/raw/main/config.example.yaml -O config.yaml
# Run http-echo-server
http-echo-server -c config.yaml
```

## Example

```json
$ curl http://127.0.0.1:80
{
  "remote_addr": "127.0.0.1:37572",
  "client_ip": "127.0.0.1",
  "http": {
    "request": {
      "method": "GET",
      "url": "/",
      "host": "127.0.0.1:80",
      "remote_addr": "127.0.0.1:37572",
      "request_uri": "/",
      "proto": "HTTP/1.1",
      "proto_major": 1,
      "proto_minor": 1,
      "header": {
        "Accept": [
          "*/*"
        ],
        "User-Agent": [
          "curl/8.2.1"
        ]
      },
      "content_length": 0,
      "close": false,
      "raw_body": "",
      "body_sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "body": ""
    }
  }
}
```
