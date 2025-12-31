<template>
  <main class="login-container">
    <h1 class="login-title">Benvenuto in Project-WASA!</h1>
    <h2 class="login-subtitle">
      Effettua il login o registrati scegliendo un nome utente per
      iniziare!
    </h2>
    <input
      v-model="nameForLogin"
      class="login-input"
      type="text"
      placeholder="Inserisci il tuo nome utente..."
    >
    <button class="login-button" @click="doLogin">Accedi</button>
  </main>
</template>

<script setup>
import axios from "../services/axios.js";
import { useRouter } from "vue-router";
import { ref } from "vue";
import { userAuthId, userName, userPhoto } from "../services/auth.js";

const router = useRouter();
const nameForLogin = ref("");

const doLogin = async () => {
	try {
		const res = await axios.post("/session", { name: nameForLogin.value });
		sessionStorage.setItem("userAuthId", res.data.id);
		userAuthId.value = res.data.id;

		sessionStorage.setItem("userName", res.data.name);
		userName.value = res.data.name;

		sessionStorage.setItem("userPhoto", res.data.photo);
		userPhoto.value = res.data.photo;

		router.push("/home");
	} catch (error) {
		// TODO
	}
};
</script>

<style scoped>
.login-container {
	font-family: sans-serif;
}
.login-title {
	text-align: center;
	margin-top: 10rem;
	margin-bottom: 0;
	font-size: 3rem;
}
.login-subtitle {
	text-align: center;
	font-size: 1.5rem;
}
.login-input {
	display: block;
	width: 20rem;
	height: 2.5rem;
	margin: 0 auto;
	border: 0.2rem solid black;
	outline: none;
	font-size: 1.2rem;
}
.login-button {
	display: block;
	width: 12rem;
	height: 3rem;
	margin: 1rem auto;
	background-color: whitesmoke;
	border: 0.2rem solid black;
	cursor: pointer;
	font-size: 1.5rem;
}
</style>
