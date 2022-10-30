<script setup>
import container from "@/lib/container";
import AppLink from "@/components/AppLink.vue";
</script>

<template>
	<nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
		<div :class="container">
			<a class="navbar-brand" href="/">Stash</a>
			<ul class="navbar-nav ml-auto">
				<AppLink app="messages">Messages</AppLink>
			</ul>
		</div>
	</nav>
</template>
