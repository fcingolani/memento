## Installation

### Docker

Run the `fcingolani/memento` image hosted on [Docker Cloud](https://hub.docker.com/r/fcingolani/memento/):

```
docker run -ti -p 3000:3000 fcingolani/memento
```

## Usage

1. Add a new replay

```
> curl -i http://127.0.0.1:3000/replays -X POST -d "player_name=fcingolani&level_number=1&level_version=1&time=6500"
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Tue, 25 Jul 2017 15:11:50 GMT
Content-Length: 153

{"id":1,"playerName":"fcingolani","levelNumber":1,"levelVersion":1,"time":6500,"FileUploadTicket":"7f94ddb5-85f5-4480-950e-5dd48e9973aa","FileData":null}
```

2. Upload binary data

```
> md5sum facepalm.png # just to check it later
0ad78e9490f2fe31e6275d91cd56ceb1 *facepalm.png

> curl -i http://127.0.0.1:3000/replays/1/file -X PUT --data-binary "@facepalm.png" -H "x-file-upload-ticket: 7f94ddb5-85f5-4480-950e-5dd48e9973aa"
HTTP/1.1 100 Continue

HTTP/1.1 200 OK
Date: Tue, 25 Jul 2017 15:12:28 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```

3. Find a replay to beat

```
> curl -i http://127.0.0.1:3000/replays/_tobeat/1/1/7000
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Tue, 25 Jul 2017 15:13:00 GMT
Content-Length: 153

{"id":1,"playerName":"fcingolani","levelNumber":1,"levelVersion":1,"time":6500,"FileUploadTicket":"00000000-0000-0000-0000-000000000000","FileData":null}
```

4. Fetch binary data

```
> curl -s http://127.0.0.1:3000/replays/1/file | md5sum
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
| LISTEN_ADDRESS    | [TCP address to listen on](Server listen address) | :3000
| DEBUG             | Enable or disable debug mode | false
| MAX_UPLOAD_BYTES  | Max allowed file upload size in bytes | 1048576

## Todo

- ☑ Limit file uploads.
- ☐ Add a security middleware.
- ☐ Change "time" for "score".
- ☐ Add parameter to `_tobeat` so you can choose if you want a bigger or smaller score to beat.
- ☐ Further separate file uploads from replay metadata.
- ☐ Add MySQL support.
- ☐ Add PostgreSQL support.
- ☐ Add CORS support.

## FAQ

### Why do you use 2 requests to save a replay, one to save the model and another to upload the file?

This was made for [a friend](https://martincerdeira.itch.io/) who uses GameMaker. GM's http library doesn't support multipart requests.