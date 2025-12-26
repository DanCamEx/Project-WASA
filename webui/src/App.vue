<template>
	<div v-if="!authId">
		<RouterView name="default" />
	</div>
	<template v-if="authId">
		<header class="app-header">
			<RouterLink to="/home" style="text-decoration: none; color: black">
				<h1 class="title">Project-WASA</h1>
			</RouterLink>
			<div>
				<button class="logout-button" @click="logout">Logout</button>
				<span class="user-info"
					>Nome: {{ username }}<br />Id utente: {{ authId }}</span
				>
			</div>
		</header>
		<div class="main-container">
			<div class="left-panel">
				<RouterView name="left" />
			</div>
			<div class="right-panel">
				<RouterView name="right" />
			</div>
		</div>
	</template>
</template>

<script setup>
import { RouterLink, RouterView, useRouter } from "vue-router";
import { ref } from "vue";

const authId = ref(localStorage.getItem("authId"));
window.authId = authId;
const username = ref("");
window.username = username;

const router = useRouter();

function logout() {
	localStorage.removeItem("authId");
	window.authId.value = localStorage.getItem("authId");
	window.username.value = "";
	router.push("/login");
}
</script>

<style>
body {
	margin: 0;
	padding: 0;
	height: 100%;
}
</style>

<style scoped>
.app-header {
	height: 3.5rem;
	background-color: #f1f1f1;
	border-bottom: 0.2rem solid black;
}
.main-container {
	display: flex;
	height: calc(100vh - 3.7rem);
}
.left-panel {
	flex: 1;
	border-right: 0.2rem solid black;
	overflow: hidden;
}
.right-panel {
	flex: 1.7;
}
.title {
	position: absolute;
	left: 12%;
	margin: 0;
	top: 0.5rem;
	font-family: sans-serif;
	cursor: pointer;
}
.logout-button {
	position: absolute;
	top: 0.7rem;
	right: 1rem;
	font-size: 1.4rem;
	font-family: sans-serif;
	cursor: pointer;
	border-radius: 0;
	border-color: #000000;
	border-style: solid;
}
.user-info {
	position: absolute;
	top: 0.7rem;
	right: 7rem;
	font-family: sans-serif;
	font-weight: bold;
	font-size: 0.8rem;
}
</style>
