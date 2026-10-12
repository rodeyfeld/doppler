/** @type {import('tailwindcss').Config} */
// One build for both the pinwheel.fan landing page and the doppler blog.
export default {
  content: ['./internal/**/*.{go,templ}', './static/js/**/*.js'],
  theme: {
    extend: {
      colors: {
        // Near-black with a hint of green, one teal, and three greys for text.
        night: '#0a0c0c',
        ink: '#101414',
        line: '#1b2120',
        'line-strong': '#2a3331',
        fg: '#e6e9e7',
        dim: '#9aa39f',
        mute: '#5f6966',
        pin: '#2dd4bf',
        rust: '#e08a7a',
      },
      fontFamily: {
        sans: ['"Instrument Sans"', 'system-ui', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'ui-monospace', 'monospace'],
      },
    },
  },
  plugins: [],
}
