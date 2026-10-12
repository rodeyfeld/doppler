package shared

import (
	"os"
	"strconv"
	"time"
)

// assetStamp is the newest modification time of the built bundles, read once
// when the process starts. Rebuilding the assets restarts the app (air in dev,
// a new image in production), so a URL like /static/css/output.css?v=<stamp>
// changes whenever the file does and browsers can't keep serving a stale copy.
var assetStamp = func() string {
	var newest time.Time
	for _, path := range []string{"static/css/output.css", "static/dist/index.js", "static/dist/index.css"} {
		if info, err := os.Stat(path); err == nil && info.ModTime().After(newest) {
			newest = info.ModTime()
		}
	}
	if newest.IsZero() {
		newest = time.Now()
	}
	return strconv.FormatInt(newest.Unix(), 36)
}()

// Asset returns a static path with a cache-busting version.
func Asset(path string) string {
	return path + "?v=" + assetStamp
}
