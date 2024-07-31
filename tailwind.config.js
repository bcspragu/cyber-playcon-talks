/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html", "./**/slides.md"],
  theme: {
    extend: {},
  },
  // Needed so that we can use utilities inline and have them override the Typography plugin
  important: true,
  plugins: [
    require('@tailwindcss/typography'),
  ],
}

