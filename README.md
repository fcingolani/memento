# Memento

[![Docker Automated build](https://img.shields.io/docker/automated/fcingolani/memento.svg)](https://hub.docker.com/r/fcingolani/memento/)

Memento is a high-score server. It can also save your replay (or other binary) data.

## Installation

### Docker

Run the `fcingolani/memento` image hosted on [Docker Hub](https://hub.docker.com/r/fcingolani/memento/):

```
docker run -ti -p 3000:3000 fcingolani/memento
```

## Usage

1. Add a new score

```
> curl -i http://127.0.0.1:3000/scores -X POST -d "player_name=fcingolani&level_number=1&level_version=1&value=6500"
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Wed, 26 Jul 2017 03:12:15 GMT
Content-Length: 177

{
  "id": 1,
  "playerName": "fcingolani",
  "levelNumber": 1,
  "levelVersion": 1,
  "value": 6500,
  "file": {
    "uploadTicket": "ce9584ad-0778-49f5-82a1-1f646121c967"
  }
}
```

2. Upload binary data

```
> md5sum facepalm.png # just to check it later
0ad78e9490f2fe31e6275d91cd56ceb1 *facepalm.png

> curl -i http://127.0.0.1:3000/files/1/data -X PUT --data-binary "@facepalm.png" -H "x-file-upload-ticket: ce9584ad-0778-49f5-82a1-1f646121c967"
HTTP/1.1 100 Continue

HTTP/1.1 200 OK
Date: Wed, 26 Jul 2017 03:50:55 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```

3. Find a score to beat

```
> curl -i "http://127.0.0.1:3000/scores/_beatable?level_number=1&level_version=1&value=7000&type=lower"
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Wed, 26 Jul 2017 03:28:01 GMT
Content-Length: 101

{
  "id": 1,
  "playerName": "fcingolani",
  "levelNumber": 1,
  "levelVersion": 1,
  "value": 6500
}
```

4. Fetch binary data

```
> curl -s http://127.0.0.1:3000/files/1/data | md5sum
0ad78e9490f2fe31e6275d91cd56ceb1 *-
```

A `/check` utility endpoint is provided, it'll respond with 200 OK so you can use it for health checks.

```
> curl -i http://localhost:3000/check
HTTP/1.1 200 OK
Date: Tue, 25 Jul 2017 20:40:58 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```

## Configuration

You can use the following environment variables to configure the server:

| Variable          | Description | Default Value
|-                  |-            |-
| DATABASE_PATH     | SQLite database path | ./db.sqlite
| LISTEN_ADDRESS    | [TCP address to listen on](https://golang.org/pkg/net/http/#Server) | :3000
| DEBUG             | Enable or disable debug mode | false
| MAX_UPLOAD_BYTES  | Max allowed file upload size in bytes | 1048576

## FAQ

### Why do you use 2 requests to save a replay, one to save the model and another to upload the file?

This was made for [a friend](https://martincerdeira.itch.io/) who uses GameMaker. GM's http library doesn't support multipart requests.

## Todo

_Not in priority order_

- ☑ Limit file uploads.
- ☐ Add a security middleware.
- ☑ Change "time" for "value".
- ☑ Add parameter to `_beatable` so you can choose if you want a bigger or smaller score to beat.
- ☑ Further separate file uploads endpoints from score metadata.
- ☐ Add MySQL support.
- ☐ Add PostgreSQL support.
- ☐ Add CORS support.
- ☐ Add Swagger file.
- ☐ Add parameter to `_beatable` to filter scores that have a file.