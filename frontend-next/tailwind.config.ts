import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./src/**/*.{js,ts,jsx,tsx,mdx}'],
  theme: {
    extend: {
      colors: {
        primary: '#0f766e',
        accent: '#155e75',
        surface: '#f8fafc',
        danger: '#b91c1c',
        success: '#166534',
        warning: '#b45309'
      }
    }
  },
  plugins: [],
};

export default config;
