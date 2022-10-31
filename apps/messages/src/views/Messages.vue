<script setup>
import { computed, onMounted, watch } from "vue";
import { useRoute } from "vue-router";
import ThreadMessages from "@/components/ThreadMessages.vue";
import ThreadSummary from "@/components/ThreadSummary.vue";
import Threads from "@/models/Threads";
import router from "@/router";

const { result, loading, error } = Threads();

const route = useRoute();

// const body = function (body) {
// 	return `<html><body>${body}</body></html>`;
// };

const current = function (id) {
	return id == route.params.id;
};

const id = computed(() => {
	return route.params.id;
});

const klass = function (id) {
	return current(id) ? "active" : "";
};

// const resize = function (e) {
// 	const iframe = e.target;
// 	iframe.height = iframe.contentWindow.document.body.scrollHeight + 40;
// };

const scroll = function (thread, top) {
	if (!thread) return;
	if (visible(thread)) return;
	document.getElementById(thread.id).scrollIntoView(top);
};

const scroll_if_current = function (thread) {
	if (current(thread.id)) {
		scroll(thread);
	}
};

const select = (thread) => {
	if (!thread) return;
	router.push({ params: { id: thread.id } });
	return thread;
};

const thread = computed(() => {
	return threads.value.find((thread) => current(thread.id));
});

const threads = computed(() => {
	return result.value?.threads.filter((t) => t);
});

const visible = function (thread) {
	const brect = document.getElementById("threads").getBoundingClientRect();
	const trect = document.getElementById(thread.id).getBoundingClientRect();
	if (trect.top < brect.top) return false;
	if (trect.top + trect.height > brect.top + brect.height) return false;
	return true;
};

onMounted(() => {
	window.addEventListener("keypress", (e) => {
		switch (e.key) {
			case "j":
				scroll(select(threads.value[threads.value.indexOf(thread.value) + 1]), false);
				break;
			case "k":
				scroll(select(threads.value[threads.value.indexOf(thread.value) - 1]), true);
				break;
		}
	});
});

watch(error, () => {
	for (const e of error.value.graphQLErrors) {
		console.log(e.message + "\n" + e.extensions?.stacktrace.join("\n"));
	}
});
</script>

<template>
	<div v-if="loading">loading</div>
	<div v-else class="row gx-0 gy-0" @keyup="key">
		<div class="col-5" id="threads">
			<ul class="list-group list-group-flush">
				<ThreadSummary
					@click="select(thread)"
					v-for="thread in threads"
					:key="thread.id"
					:thread="thread"
					:class="klass(thread.id)"
					@mounted="scroll_if_current(thread)"
				/>
			</ul>
		</div>
		<div class="col-7">
			<ThreadMessages v-for="thread in threads" :key="thread.id" :thread="thread" :current="id" />
		</div>
	</div>
</template>
