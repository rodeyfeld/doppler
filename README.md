# doppler
Personal website for rodeyfeld.

## local dev 

- Create and Migrate db using `goose`:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
touch /app/doppler.db
goose -dir /app/internal/db/sql/ sqlite3 /app/doppler.db up
```

- To install dependencies:

```bash
bun install
```

- To run Go hot reload:
```bash
air
```


To watch templ file changes:
```bash
templ generate --watch
```


To regenerate Tailwind output:
```bash
bun run tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
```


## docker

Build dockerfile
```bash
docker build --tag doppler . 
docker run -p 1323:1323 doppler
```
