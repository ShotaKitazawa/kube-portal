/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx}',
    './src/components/**/*.{js,ts,jsx,tsx}',
  ],
  darkMode: 'media', // 'media' or 'class'
  theme: {
    extend: {
      // Add custom theme settings here
    },
  },
  // variants section is not needed in Tailwind CSS v4
  plugins: [],
  // Tailwind CSS v4 new features
  future: {
    // Enable future features
    hoverOnlyWhenSupported: true,
  },
}
