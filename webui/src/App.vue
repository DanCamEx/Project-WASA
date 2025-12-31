<template>
  <section v-if="!userAuthId"><RouterView name="default" /></section>
  <template v-if="userAuthId">
    <header class="app-header">
      <div class="header-content">
        <RouterLink
          to="/home"
          style="text-decoration: none; color: black"
        >
          <h1 class="title">Project-WASA</h1>
        </RouterLink>
        <div class="header-actions">
          <button
            class="settings-button"
            @click="router.push('/settings')"
          >
            Impostazioni
          </button><span class="user-info">Nome: {{ userName }}<br>Id utente:
            {{ userAuthId }}</span><img
            class="profile-photo"
            :src="
              isImageSet(userPhoto)
                ? serverUrl + '/images/' + userPhoto
                : '/images/default_user_photo.jpg'
            "
            alt="Photo"
          ><button class="logout-button" @click="logout">
            Logout
          </button>
        </div>
      </div>
    </header>
    <main class="main-container">
      <div class="left-panel"><RouterView name="left" /></div>
      <div class="right-panel"><RouterView name="right" /></div>
    </main>
  </template>
</template>

<script setup>
import { RouterLink, RouterView, useRouter } from "vue-router";
import { userAuthId, userName, userPhoto } from "./services/auth.js";
import { isImageSet } from "./services/utilities.js";

const router = useRouter();
const serverUrl = __API_URL__;

const logout = () => {
	sessionStorage.clear();
	userAuthId.value = null;
	userName.value = null;
	userPhoto.value = null;
	router.push("/login");
};
</script>

<style>
body {
	margin: 0;
	padding: 0;
	font-family: sans-serif;
}
</style>

<style scoped>
.app-header {
	display: flex;
	align-items: center;
	height: 3.5rem;
	background-color: whitesmoke;
	border-bottom: 0.2rem solid black;
}
.header-content {
	display: flex;
	align-items: center;
	justify-content: space-between;
	width: 100%;
	padding: 0 2rem;
}
.header-actions {
	display: flex;
	align-items: center;
	gap: 1.2rem;
}
.main-container {
	display: flex;
	height: calc(100vh - 3.7rem);
}
.left-panel {
	flex: 1;
	border-right: 0.2rem solid black;
	overflow-y: auto;
}
.right-panel {
	flex: 1.7;
	overflow-y: auto;
}
.title {
	margin: 0;
	cursor: pointer;
	font-family: sans-serif;
}
.settings-button,
.logout-button {
	cursor: pointer;
	padding: 0.3rem 0.7rem;
	border-radius: 0;
	border-color: black;
	border-style: solid;
	font-size: 1rem;
	font-family: sans-serif;
}
.user-info {
	font-size: 0.8rem;
	font-family: sans-serif;
	font-weight: bold;
}
.profile-photo {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	border: 2px solid black;
	object-fit: cover;
	object-position: center;
}
</style>
