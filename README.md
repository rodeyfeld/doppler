# doppler

To install dependencies:

```bash
bun install
```


To run Go hot reload:
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


Run all
```bash
air & templ generate --watch & bun run tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
```
