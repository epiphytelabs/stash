<script setup>
import { computed } from "vue";

const props = defineProps(["address"]);

const name = computed(() => props.address.name);
const address = computed(() => props.address.address);
</script>

<template>
	<div v-if="name" class="text-truncate">
		{{ name }}
	</div>
	<div v-else-if="address" class="text-truncate">
		{{ address }}
	</div>
</template>
