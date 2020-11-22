# Snippets

Interview take-home test question for Stellar.

## Assumptions

- This is running behind a tls proxy, so that's not necessary to handle in this.
- The instructions mean that `expires_in` is in seconds.
- Since the only used `Content-Type` is `application/json`, it's fine to accept any and just assume it is json. We don't need strict enforcement.

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

## Todo

- [x] store snippets
- [ ] retrieve snippets
- [ ] expire snippets
- [ ] reset expiry when retrieving snippets
- [ ] Add tests
- [ ] Either handle tls, or put this behind a tls proxy.
- [ ] Structured logging
- [ ] Return an HTTP 413 error if body is longer than 1mb.
- [ ] authentication
- [ ] rate-limiting
