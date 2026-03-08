/** @type {import('next').NextConfig} */
const nextConfig = {
  basePath: '/app',
  poweredByHeader: false,
  async redirects() {
    return [
      {
        source: '/',
        destination: '/app',
        permanent: false,
        basePath: false,
      },
    ];
  },
};

export default nextConfig;
