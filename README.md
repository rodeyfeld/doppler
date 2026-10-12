# Doppler

My personal photoblog thing. Live at https://pinwheel.fan

## Running it

   ```bash
   cp example.env .env
   docker compose up
   ```

Then hit http://localhost:1323

## Deploying

Build and push:
```bash
docker build -t edrodefeld/doppler .
docker push edrodefeld/doppler:latest
```

Deploy to k8s:
```bash
kubectl apply -f ../mirage/deployments/doppler.yml
kubectl rollout restart deployment/doppler -n galaxy
```

## What's in it

Go + Echo for the backend, Templ for templates, Tailwind for styling. HTMX makes it interactive, Quill for the rich text editor, and tsParticles for that space background effect. Photos go in S3 (Garage), metadata in SQLite.

The look is one small design system for both the pinwheel.fan landing page and the blog: plain Tailwind with its own tokens in `tailwind.config.js` and the component classes in `static/css/input.css` (no daisyUI). Fonts are Instrument Sans and JetBrains Mono from Google Fonts. The landing page's project list lives in `internal/components/home/home.templ`.

Air watches for changes and automatically rebuilds everything - runs `templ generate` for templates, `bun run build:js` for JavaScript, and `bun run build:css` for Tailwind. Then restarts the Go app. Air doesn't watch `.css` files, so after editing `input.css` run `docker compose exec web bun run build:css` (or touch a template).

## Sample data

With the compose stack up, fill the blog with sample posts (and generated photos in Garage):

```bash
docker compose exec web go run ./cmd/seed          # refuses if posts already exist
docker compose exec web go run ./cmd/seed -reset   # wipes ALL posts first
```

Posts are by a sample user, `pinwheel` (password `pinwheel`), created if missing.
