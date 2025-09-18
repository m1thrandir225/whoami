import { defineConfig } from 'vite'
import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

import { TanStackRouterVite } from '@tanstack/router-plugin/vite'
import { resolve } from 'node:path'
import fs from 'node:fs'
import path from 'node:path'

const certDir = process.env.VITE_CERT_DIR
	? process.env.VITE_CERT_DIR
	: path.resolve(__dirname, '../deployment/certs')

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [
		TanStackRouterVite({ autoCodeSplitting: true }),
		viteReact(),
		tailwindcss(),
	],
	test: {
		globals: true,
		environment: 'jsdom',
	},
	server: {
		port: 5173,
		https: {
			key: fs.readFileSync(path.join(certDir, 'localhost-key.pem')),
			cert: fs.readFileSync(path.join(certDir, 'localhost-cert.pem')),
		},
	},

	resolve: {
		alias: {
			'@': resolve(__dirname, './src'),
		},
	},
})
