module.exports = {
  content: [
    './templates/**/*.{html,templ}',
    './internal/**/*.{go,templ}',
    './static/css/tailwind-safelist.txt'
  ],
  theme: {
    extend: {},
  },
  plugins: [],
  safelist: [
    'space-x-8',
    'hidden',
    'md:flex',
    'text-gray-600',
    'hover:text-gray-900',
    'font-medium'
  ]
}