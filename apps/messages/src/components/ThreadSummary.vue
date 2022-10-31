<script setup>
import { computed, onMounted } from "vue";
import RelativeTime from "@/components/RelativeTime.vue";

const emit = defineEmits(["mounted"]);

const props = defineProps(["thread"]);

const participants = computed(() => {
	return props.thread.messages.reduce((ax, message) => {
		if (!message.from) return ax;
		if (!ax.includes(message.from.display)) {
			ax.push(message.from.display);
		}
		return ax;
	}, []);
});

const subject = computed(() => {
	return props.thread.messages[0].subject;
});

onMounted(() => {
	emit("mounted", this);
});
</script>

<template>
	<li class="list-group-item thread" :id="props.thread.id">
		<div class="d-flex mb-1">
			<div class="flex-grow-1 participants text-truncate">
				<span class="me-2" v-for="participant in participants" :key="participant">
					{{ participant }}
				</span>
			</div>
			<RelativeTime class="flex-shrink-0 ms-2 updated" :time="props.thread.updated" />
		</div>
		<div class="subject text-truncate">{{ subject }}</div>
	</li>
</template>
