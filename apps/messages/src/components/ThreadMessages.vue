<script setup>
import { computed } from "vue";
import MessageAddress from "@/components/MessageAddress.vue";
import MessageBody from "@/components/MessageBody.vue";
import RelativeTime from "@/components/RelativeTime.vue";

const props = defineProps(["current", "thread"]);

const klass = computed(() => {
	return props.current == props.thread.id ? "visible" : "hidden";
});

const thread = computed(() => {
	return props.thread;
});
</script>

<template>
	<div class="messages p-3 pb-0" :class="klass">
		<div class="subject">{{ thread.subject }}</div>
		<div class="card mb-3 message" v-for="message in thread.messages" :key="message.id">
			<div class="card-header d-flex">
				<MessageAddress :address="message.from" class="flex-grow-1" />
				<RelativeTime :time="message.received" class="flex-shrink-0" />
			</div>
			<div class="card-body p-0">
				<MessageBody :body="message.body" />
			</div>
		</div>
	</div>
</template>
