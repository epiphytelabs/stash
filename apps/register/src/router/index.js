import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	linkActiveClass: "active",
	linkExactActiveClass: "active",
	routes: [
		{
			path: "",
			name: "registration",
			component: () => import("../views/Registration.vue"),
		},
	],
});

export default router;
