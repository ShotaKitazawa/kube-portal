/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'export',
  images: {
    unoptimized: true,
  },
  // Tailwind CSS v4との互換性を確保
  webpack: (config) => {
    return config
  },
}

module.exports = nextConfig
