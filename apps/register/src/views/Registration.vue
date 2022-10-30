<script setup>
import { ref, onMounted } from "vue";
import axios from "axios";

const first = ref(null);
const id = ref(null);
const password = ref(null);

onMounted(() => {
	first.value.focus();
});

const register = async () => {
	const data = new FormData();
	data.append("id", id.value);
	data.append("password", password.value);

	const res = await axios.post("/api/users", data);

	const link = document.createElement("a");
	link.href = `data:text/plain;base64,${res.data}`;
	link.download = `${id.value}.p12`;
	link.click();
};
</script>

<template>
	<div class="card col-4 offset-4">
		<div class="card-header">Account Registration</div>
		<div class="card-body">
			<form>
				<div class="mb-3">
					<label for="id" class="form-label">Email</label>
					<input ref="first" v-model="id" type="email" class="form-control" id="id" />
				</div>
				<div class="mb-3">
					<label for="password" class="form-label">Password</label>
					<input v-model="password" type="password" class="form-control" id="password" />
				</div>
				<button class="btn btn-primary align-end" @click.prevent="register()">Register</button>
			</form>
		</div>
	</div>
</template>
