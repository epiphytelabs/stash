import { createApp, h } from "vue";
import App from "./App.vue";
const app = createApp({
	render: () => h(App),
});

import "bootstrap";

import { createPinia } from "pinia";
app.use(createPinia());

import router from "./router";
app.use(router);

app.mount("#app");
