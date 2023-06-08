# nomad-exercise

## Description

This project is a simple API that serves as an exercise to explore the Nomad API.

## API Endpoint

The API provides the following endpoint:

### POST /services

This endpoint starts a Nomad job and accepts the following parameters:

- `name` (string): A unique identifier for the service to create.
- `url` (string): The URL of the page to serve.
- `script` (boolean): Specifies whether the `url` should be interpreted as a script or a page.

The endpoint performs the following actions based on the provided parameters:

- If `script` is set to `false`, the `url` is downloaded and served as `index.html`.
- If `script` is set to `true`, the `url` is downloaded, the script is executed, and its output is downloaded as `index.html`.

The endpoint returns the URL of the page.

## Example Usage

### Non-Script Example

```bash
curl localhost:3000/services -X POST -d '{"name": "my-page", "url": "https://example.com/mypage.html", "script": false}'
```

Response:

```json
{
    "url": "http://192.168.1.104:27510"
}
```

Assuming `https://example.com/mypage.html` contains the following content:

```
Hello World!
```

Then the output of `curl http://192.168.1.104:27510` will be:

```
Hello World!
```

### Script Example

```bash
curl localhost:3000/services -X POST -d '{"name": "my-page", "url": "https://example.com/myscript.sh", "script": true}'
```

Response:

```json
{
    "url": "http://192.168.1.104:27511"
}
```

Assuming `https://example.com/myscript.sh` contains the following content:

```bash
#!/bin/sh

echo Hello World!
```

Then the output of `curl http://192.168.1.104:27511` will be:

```
Hello World!
```

## Installation

To install and run this project locally, follow these steps:

1. Start Nomad by running the command: `nomad agent -dev -bind 0.0.0.0 -network-interface='{{ GetDefaultInterfaces | attr "name" }}'`.
2. Start the API by running the command: `go run ./cmd/api`.
