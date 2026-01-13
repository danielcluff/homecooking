// @ts-check
import { defineConfig } from 'astro/config';
import tailwindcss from '@tailwindcss/vite';
import typography from '@tailwindcss/typography';

// https://astro.build/config
export default defineConfig({
  integrations: [
    tailwindcss({
      plugins: [typography()],
    }),
  ],
  output: 'static',
  server: {
    port: 3001,
  },
});
