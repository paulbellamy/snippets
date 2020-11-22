# Snippets

Interview take-home test question for Stellar.

## Dependencies

- docker-compose

## Running

```
$ docker-compose up
```

Then, once that is running visit http://localhost:3000.

## Usage

```sh
curl -X POST -H "Content-Type: application/json" -d '{"name":"recipe", "expires_in": 30, "snippet":"1 apple"}' https://localhost:3000/snippets
# response 201 Created
{
  "url": "https://localhost:3000/snippets/recipe",
  "name": "recipe",
  "expires_at": "2020-02-22T20:02:02Z",
  "snippet": "1 apple"
}

curl https://localhost:3000/snippets/recipe
# response 200 OK
{
  "url": "https://localhost:3000/snippets/recipe",
  "name": "recipe",
  "expires_at": "2020-02-22T20:02:32Z",
  "snippet": "1 apple"
}

# wait 60 seconds
curl https://example.com/snippets/recipe
# response 404 Not Found
```

## Todo

- [ ] store snippets
- [ ] retrieve snippets
- [ ] expire snippets
- [ ] Add tests
- [ ] Either handle tls, or put this behind a tls proxy.
- [ ] Structured logging
- [ ] Return an HTTP 413 error if body is longer than 1mb.
- [ ] authentication
- [ ] rate-limiting
