<script setup>
import { RouterView } from "vue-router";
import NavigationBar from "@/components/NavigationBar.vue";
import container from "@/lib/container";
</script>

<template>
	<header><NavigationBar /></header>

	<div :class="container" style="margin-top: 72px">
		<RouterView />
	</div>
</template>

<style lang="scss">
@import "@/styles/bootstrap.scss";
@import "bootstrap";
@import "@/styles/app.scss";
</style>
