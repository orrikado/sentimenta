import { paraglideVitePlugin } from '@inlang/paraglide-js';
import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:8000', // Local backend
				changeOrigin: true
			},
			'/ws': {
				target: 'ws://localhost:8000',
				ws: true
			}
		}
	},
	build: {
		minify: 'esbuild',
		cssMinify: 'esbuild'
	},
	plugins: [
		tailwindcss(),
		sveltekit(),
		paraglideVitePlugin({
			project: './project.inlang',
			outdir: './src/lib/paraglide'
		})
	]
});
