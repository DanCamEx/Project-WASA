<template>
  <template v-if="!userAuthId"><RouterView name="default" /></template>
  <template v-if="userAuthId">
    <header class="main-header">
      <RouterLink
        to="/home"
        style="
          text-decoration: none;
          color: black;
          align-self: center;
          margin-left: 0.5rem;
        "
      >
        <h1 style="margin: 0; margin-left: 0.5rem">Project-WASA</h1>
      </RouterLink>
      <div class="right-header-side">
        <button class="settings-button" @click="router.push('/settings')">
          Impostazioni
        </button>
        <p class="user-info">
          Nome: {{ userName }}<br>Id utente: {{ userAuthId }}
        </p>
        <img
          class="profile-photo"
          :src="
            isImageSet(userPhoto)
              ? serverUrl + '/images/' + userPhoto
              : '/images/default_user_photo.jpg'
          "
          alt="Photo"
        ><button class="logout-button" @click="logout">Logout</button>
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
  font-family: "Helvetica Neue", Arial, sans-serif;
}
button,
input,
select,
textarea {
  font-family: inherit;
}
</style>

<style scoped>
.main-header {
  display: flex;
  min-height: 3.5rem;
  background-color: gainsboro;
  border-bottom: 0.2rem solid black;
}
.right-header-side {
  display: flex;
  align-items: center;
  margin-left: auto;
  gap: 0.5rem;
}
.settings-button,
.logout-button {
  align-self: center;
  margin-left: auto;
  margin-right: 1rem;
  cursor: pointer;
  padding: 0.4rem;
  border: 0.15rem solid black;
  background-color: rgb(0, 150, 250);
  font-size: 1rem;
  font-weight: bold;
}
.settings-button:hover,
.logout-button:hover {
  background-color: rgb(0, 100, 200);
}
.user-info {
  margin: 0 0 0 0.5rem;
  align-self: center;
  font-size: 0.8rem;
  font-weight: bold;
  white-space: nowrap;
}
.profile-photo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 2px solid black;
  object-fit: cover;
  object-position: center;
}
.main-container {
  display: flex;
  height: calc(100vh - 3.7rem);
}
.left-panel {
  flex: 1;
  border-right: 0.2rem solid black;
  overflow: auto;
}
.right-panel {
  flex: 1.7;
  overflow: auto;
}
</style>
