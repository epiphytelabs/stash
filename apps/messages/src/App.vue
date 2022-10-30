<script setup>
import { RouterView } from "vue-router";
import Navbar from "@/components/Navbar.vue";
</script>

<template>
	<header><Navbar /></header>

	<div style="margin-top: 56px">
		<RouterView />
	</div>
</template>

<style lang="scss">
@import "@/styles/bootstrap.scss";
@import "bootstrap";
@import "@/styles/app.scss";
@import "vue-select/dist/vue-select.css";
</style>
