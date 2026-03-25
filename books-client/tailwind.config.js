/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        bg: {
          main: '#0d0e1c',
          dim: '#1d2235',
          active: '#4a4f69',
          inactive: '#2b3045',
        },
        fg: {
          main: '#ffffff',
          dim: '#989898',
          alt: '#c6daff',
        },
        border: '#61647a',
        red: {
          DEFAULT: '#ff5f59',
          faint: '#ef8386',
        },
        green: {
          DEFAULT: '#44bc44',
          faint: '#88ca9f',
        },
        yellow: {
          DEFAULT: '#d0bc00',
          faint: '#d2b580',
        },
        blue: {
          DEFAULT: '#2fafff',
          faint: '#82b0ec',
        },
        magenta: {
          DEFAULT: '#feacd0',
          faint: '#caa6df',
        },
        cyan: {
          DEFAULT: '#00d3d0',
          faint: '#9ac8e0',
        },
      },
    },
  },
  plugins: [],
}
