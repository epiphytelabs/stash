import { fileURLToPath, URL } from "url";
import { defineConfig } from "vite";

// plugins
import vue from "@vitejs/plugin-vue";
import viteGraphQLPlugin from "vite2-graphql-plugin";
import pluginRewriteAll from "vite-plugin-rewrite-all";

export default defineConfig({
	base: "/apps/messages/",
	build: {
		outDir: "./dist/web",
	},
	plugins: [vue(), viteGraphQLPlugin(), pluginRewriteAll()],
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
