/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./internal/**/*.{go,templ,js,html}'],
  theme: {
    extend: {},
  },
  plugins: [require('daisyui')],
  daisyui: {
    themes: ["forest"],
  }
}

