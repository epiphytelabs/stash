import { fileURLToPath, URL } from "url";
import { defineConfig } from "vite";

// plugins
import vue from "@vitejs/plugin-vue";
import viteGraphQLPlugin from "vite2-graphql-plugin";

export default defineConfig({
	plugins: [vue(), viteGraphQLPlugin()],
	resolve: {
		alias: {
			"@": fileURLToPath(new URL("./src", import.meta.url)),
		},
	},
	server: {
		hmr: {
			clientPort: process.env.VITE_CLIENT_PORT,
		},
	},
});
