# Snippets

Simple rest server to store expiring bits of text.

## Dependencies

- docker-compose

## Running

```
$ docker-compose up
```

Then, once that is running visit http://localhost:3000.

## Usage

```sh
curl -X POST -H "Content-Type: application/json" -d '{"name":"recipe", "expires_in": 30, "snippet":"1 apple"}' http://localhost:3000/snippets
# response 201 Created
{
  "url": "http://localhost:3000/snippets/recipe",
  "name": "recipe",
  "expires_at": "2020-02-22T20:02:02Z",
  "snippet": "1 apple"
}

curl http://localhost:3000/snippets/recipe
# response 200 OK
{
  "url": "http://localhost:3000/snippets/recipe",
  "name": "recipe",
  "expires_at": "2020-02-22T20:02:32Z",
  "snippet": "1 apple"
}

# wait 60 seconds
curl http://localhost:3000/snippets/recipe
# response 404 Not Found
```

## Assumptions

- This is running behind a tls proxy, so that's not necessary to handle in this.
- The instructions mean that `expires_in` is in seconds.
- Since the only used `Content-Type` is `application/json`, it's fine to accept any and just assume it is json. We don't need strict enforcement.
- Multiple POSTS with the same name should override the existing snippet (as
  opposed to erroring).

## Design Decisions

- For persistence I've chosen Redis. No point in reinventing our own persistence
    format for this exercise. I'm familiar with redis, but more importantly, it
    has built-in expiry for keys. So, implementing this with it is really
    straightforward, and doesn't need our own garbage-collector implementation.
    - I've left in the in-memory implementation for reference, or perhaps unit
        tests.
        To see that you could run:
        ```sh
        go run . -db "memory://"
        ```
- Error handling is pretty basic. Where appropriate I've returned it to the user
    in the http response. Otherwise I've logged and masked it with HTTP 500, to
    avoid leaking implementation details to the user.
- There's a bunch of other stuff to do, to get it production ready (see below
    Todo list), but some of it is dependent on the deployment environment, (e.g.
    TLS, auth, rate-limiting), so I've left that off (also time-limits).
- It would be great to have unit tests, but I've left them off as well due to
    time.

## Todo

- [x] store snippets
- [x] retrieve snippets
- [x] expire snippets
- [x] reset expiry when retrieving snippets
- [x] logging
- [x] persist snippets on disk
- [ ] Return an HTTP 413 error if request body is too long.
- [ ] Add tests
- [ ] Either handle tls, or put this behind a tls proxy.
- [ ] authentication & rate-limiting
