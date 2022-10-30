<script setup>
import { computed } from "vue";

const props = defineProps(["body"]);

const body = computed(() => {
	switch (true) {
		case props.body.html != "":
			return props.body.html.replace(/<p[^>]*>—(.*?)<\/p>/gm, "");
		case props.body.text != "":
			return "<p>" + props.body.text.replace(/\n/g, "<br>") + "</p>";
		default:
			return "";
	}
});
</script>

<template>
	<div class="body" v-html="body" />
</template>
