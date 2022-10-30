<script setup>
import { onMounted } from "vue";
import RelativeTime from "@/components/RelativeTime.vue";

const emit = defineEmits(["mounted"]);

const props = defineProps(["thread"]);

onMounted(() => {
	emit("mounted", this);
});
</script>

<template>
	<li class="list-group-item thread" :id="props.thread.id">
		<div class="d-flex mb-1">
			<div class="flex-grow-1 participants text-truncate">
				<span class="me-2" v-for="participant in props.thread.participants" :key="participant.display">
					{{ participant.display }}
				</span>
			</div>
			<RelativeTime class="flex-shrink-0 ms-2 updated" :time="props.thread.updated" />
		</div>
		<div class="subject text-truncate">{{ props.thread.subject }}</div>
	</li>
</template>
