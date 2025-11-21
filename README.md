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

Air watches for changes and automatically rebuilds everything - runs `templ generate` for templates, `bun run build:js` for JavaScript, and `bun run build:css` for Tailwind. Then restarts the Go app.
