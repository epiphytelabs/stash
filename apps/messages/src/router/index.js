import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	linkActiveClass: "active",
	linkExactActiveClass: "active",
	routes: [
		{
			path: "/:id(.+)?",
			name: "messages",
			component: () => import("../views/Messages.vue"),
		},
	],
});

export default router;
