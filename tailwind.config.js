/** @type {import('tailwindcss').Config} */
import daisyui from "daisyui"
export default {
  content: ['./internal/**/*.{go,templ,js,html}'],
  theme: {
    extend: {},
  },
  plugins: [daisyui,],
  daisyui: {
    themes: ["forest"],
  }
}

