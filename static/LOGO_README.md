# Doppler Logo Variations

The Doppler logo features a stylized pinwheel that represents the concept of "emitting signals" - posts radiating outward like broadcast transmissions.

## Logo Files

### `doppler-logo-header.svg` (Header/Navbar)
Compact animated version specifically designed for the header navbar:
- Moderately sized pinwheel (65% scale)
- Ultra-compact signal waves (max radius 125px)
- Subtle pulsing center
- Spin triggered only on hover/focus via CSS
- Fits perfectly in small navbar spaces
- **Use for:** Header navigation, tight spaces where logo needs to be contained

### `doppler-logo.svg` (Full Animated)
The full animated logo featuring:
- Spinning pinwheel with gradient colors (green → cyan → purple)
- Pulsing signal waves emanating outward from the center
- Central pulse effect representing the signal source
- **Use for:** Hero sections, large branding areas, promotional materials

### `doppler-logo-simple.svg` (Static)
A simplified static version featuring:
- Static pinwheel with the same gradient
- Fixed signal wave rings
- Subtle central glow
- **Use for:** Print materials, contexts where animation isn't appropriate

### `pinwheel.svg` (Legacy)
The original plain pinwheel without signal elements.
- Kept for backward compatibility

## Design Concept

The logo combines:
1. **Pinwheel**: Represents the Pinwheel Labs brand identity
2. **Signal Waves**: Concentric circles radiating outward symbolize broadcast/transmission
3. **Signal Dots**: Four dots at cardinal positions represent individual posts/signals
4. **Rotation**: The spinning pinwheel creates a dynamic, living feel
5. **Colors**: Green-cyan-purple gradient matches the "forest" theme

## Technical Details

- **Format**: SVG (scalable vector graphics)
- **Dimensions**: 512×512px viewBox (scales to any size)
- **Animation**: CSS animations (compatible with all modern browsers)
- **Colors**:
  - Primary green: `#4ade80`
  - Cyan: `#22d3ee`
  - Purple: `#a855f7`
  - Dark stroke: `#1f2937`

## Usage in Templates

```go
// Header/navbar (compact, fits in small spaces)
<a class="doppler-brand inline-flex items-center gap-3 text-xl group" href="/doppler/">
  <span class="relative flex items-center justify-center w-12 h-12 overflow-visible">
    <span class="doppler-logo-backdrop absolute inset-0 flex items-center justify-center" aria-hidden="true">
      <i class="fa-solid fa-record-vinyl text-primary/35 text-3xl"></i>
    </span>
    <img src="/static/doppler-logo-header.svg"
         alt="Doppler"
         class="doppler-logo-img relative z-10 w-full h-full pointer-events-none select-none" />
  </span>
  <span class="font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">doppler</span>
</a>

// Full animated version (for larger spaces)
<img src="/static/doppler-logo.svg" alt="Doppler" class="w-16 h-16" />

// Static version (for print/performance-critical contexts)
<img src="/static/doppler-logo-simple.svg" alt="Doppler" class="w-10 h-10" />
```

**Note**: Each logo variant is optimized for its specific use case:
- **Header logo** (`doppler-logo-header.svg`): Compact design fits perfectly in navbar buttons and tight spaces (32px recommended). Pair it with a translucent vinyl backdrop (`fa-record-vinyl`), avoid containers that add padding/overflow clipping (e.g., DaisyUI `.btn`), keep the logo span `overflow-visible`, and trigger spin via CSS on hover/focus.
- **Full logo** (`doppler-logo.svg`): More expansive signal waves for hero sections and larger displays (64px+)
- **Static logo** (`doppler-logo-simple.svg`): Non-animated fallback for print or reduced motion preferences

## Animation Timings

- **Pinwheel rotation**: 8 seconds per full rotation
- **Signal waves**: 2-second pulse cycle (staggered at 0.5s intervals)
- **Signal dots**: 1.5-second blink cycle (staggered at 0.25s-1s intervals)
- **Central pulse**: 2-second breathing effect

