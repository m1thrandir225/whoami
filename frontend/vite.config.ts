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

const isDockerBuild =
	process.env.NODE_ENV === 'production' || process.env.DOCKER_BUILD === 'true'

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [
		TanStackRouterVite({ autoCodeSplitting: true }),
		viteReact(),
		tailwindcss(),
	],

	server: {
		port: 5173,
		// Only use HTTPS in development, not in Docker build
		...(isDockerBuild
			? {}
			: {
					https: {
						key: fs.readFileSync(path.join(certDir, 'localhost-key.pem')),
						cert: fs.readFileSync(path.join(certDir, 'localhost-cert.pem')),
					},
				}),
	},

	resolve: {
		alias: {
			'@': resolve(__dirname, './src'),
		},
	},
})
