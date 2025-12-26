<template>
	<h1 class="login-title">Benvenuto in Project-WASA!</h1>
	<h2 class="login-subtitle">
		Effettua il login o registrati scegliendo un nome utente per iniziare!
	</h2>
	<input
		type="text"
		v-model="name"
		placeholder="Inserisci il tuo nome utente..."
		class="login-input"
	/>
	<button class="login-button" @click="doLogin">Accedi</button>
</template>

<script setup>
import axios from "../services/axios.js";
import { useRouter } from "vue-router";
import { ref } from "vue";

const router = useRouter();
const name = ref("");
async function doLogin() {
	try {
		const response = await axios.post("/session", {
			name: name.value,
		});
		const id = response.data.id;
		localStorage.setItem("authId", id);
		window.authId.value = id;
		router.push("/home");
	} catch (error) {
		console.error("Errore durante il login:", error);
	}
}
</script>
<style scoped>
.login-title {
	font-size: 3rem;
	font-family: sans-serif;
	text-align: center;
	margin-top: 10rem;
	margin-bottom: 0;
}
.login-subtitle {
	font-size: 1.5rem;
	font-family: sans-serif;
	text-align: center;
}
.login-input {
	display: block;
	margin: 0rem auto;
	width: 20rem;
	height: 2.5rem;
	font-size: 1.2rem;
	font-family: sans-serif;
	border: 0.2rem solid #000000;
	outline: none;
}
.login-button {
	display: block;
	margin: 1rem auto;
	width: 12rem;
	height: 3rem;
	font-size: 1.5rem;
	font-family: sans-serif;
	border: 0.2rem solid #000000;
	background-color: #f1f1f1;
	cursor: pointer;
}
</style>
