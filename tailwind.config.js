/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "class",
  content: [
    "./views/**/*.html",
    "./node_modules/flowbite/**/*.js",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          ...require("tailwindcss/colors").blue,
        }
      }
    },
  },
  plugins: [
    require("flowbite/plugin"),
  ],
}

