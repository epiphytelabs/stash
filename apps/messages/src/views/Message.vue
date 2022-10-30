<script setup>
import Message from "@/models/Message";
import MessageAddress from "@/components/MessageAddress.vue";
import RelativeTime from "@/components/RelativeTime.vue";
import { useRoute } from "vue-router";

const route = useRoute();

const { result, error } = Message(route.params.hash);
</script>

<template>
	<div class="pt-3">
		<div v-if="error" class="alert alert-danger">{{ error }}</div>
		<div v-else-if="loading"></div>
		<div v-else>
			<h4 class="mb-3 fw-bold">{{ result.message.subject }}</h4>
			<div class="card">
				<div class="card-header d-flex">
					<div class="flex-shrink-0">
						<MessageAddress :address="result.message.from" />
					</div>
					<div class="flex-grow-1 text-end">
						<RelativeTime :time="result.message.received" />
					</div>
				</div>
				<iframe seamless sandbox :srcdoc="result.message.body.html" width="100%" height="600px" />
			</div>
		</div>
	</div>
</template>
