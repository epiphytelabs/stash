<script setup>
const props = defineProps(["error"]);
</script>

<template>
	<div class="alert alert-danger" v-for="error in props.error.graphQLErrors" :key="error">
		<h6 class="mb-0">{{ error.message }}</h6>
		<div class="font-monospace mt-2" style="font-size: 0.7em" v-for="line in error.extensions?.stacktrace" :key="line">
			{{ line }}
		</div>
	</div>
</template>
