<script setup>
import { computed } from "@vue/reactivity";

const props = defineProps({
	app: String,
	active: Boolean,
});

const klass = computed(() => {
	return props.active ? "active" : "nav-link";
});

const href = `/apps/${props.app}/`;
</script>

<template>
	<li class="nav-item ms-3">
		<a :href="href" class="nav-link" :class="klass">
			<slot />
		</a>
	</li>
</template>
