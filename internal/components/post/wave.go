package post

import (
	"fmt"
	"time"
)

// waveClass picks the wavelength of the wave drawn above a post: tight ripples
// for new posts, long swells for old ones. The classes set --wl in input.css and
// live outside Tailwind's layers so the names survive purging.
func waveClass(created time.Time) string {
	days := time.Since(created).Hours() / 24
	switch {
	case days < 7:
		return "wave-1"
	case days < 21:
		return "wave-2"
	case days < 45:
		return "wave-3"
	case days < 90:
		return "wave-4"
	case days < 180:
		return "wave-5"
	default:
		return "wave-6"
	}
}

// ageLabel is the tooltip on a post's date: how long ago it was posted.
func ageLabel(created time.Time) string {
	days := int(time.Since(created).Hours() / 24)
	switch {
	case days < 1:
		return "today"
	case days == 1:
		return "yesterday"
	case days < 30:
		return fmt.Sprintf("%d days ago", days)
	case days < 60:
		return "a month ago"
	case days < 365:
		return fmt.Sprintf("%d months ago", days/30)
	case days < 730:
		return "a year ago"
	default:
		return fmt.Sprintf("%d years ago", days/365)
	}
}
