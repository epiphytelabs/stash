<script setup>
import AppLink from "@/components/AppLink.vue";
</script>

<template>
	<nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
		<div class="container-fluid">
			<a class="navbar-brand" href="/">&nbsp;Stash</a>
			<ul class="navbar-nav ml-auto">
				<AppLink app="messages" active>Messages</AppLink>
			</ul>
		</div>
	</nav>
</template>
