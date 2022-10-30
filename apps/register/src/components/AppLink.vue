<script setup>
const props = defineProps(["app"]);

const href = `https://${import.meta.env.VITE_DOMAIN}/apps/${props.app}/`;
</script>

<template>
	<li class="nav-item ms-3">
		<a :href="href" class="nav-link">
			<slot />
		</a>
	</li>
</template>
