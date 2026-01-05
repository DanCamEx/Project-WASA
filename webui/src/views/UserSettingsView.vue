<template>
  <div class="dark-overlay" />
  <main style="padding: 1rem">
    <h1>Impostazioni Utente</h1>
    <ul style="list-style-type: none; padding: 0">
      <strong style="margin-left: 3rem">Immagine</strong>
      <li
        style="
          display: flex;
          flex-direction: row;
          gap: 0.5rem;
          margin-top: 0.3rem;
        "
      >
        <img
          :src="
            isImageSet(newPhoto)
              ? serverUrl + '/images/' + newPhoto
              : '/images/default_user_photo.jpg'
          "
          alt="Foto Utente"
          width="150"
          height="150"
          class="user-photo"
        >
        <input
          id="photoInput"
          type="file"
          :hidden="true"
          accept="/*"
          @click="$event.target.value = null"
          @change="updateUserPhoto($event.target.files[0])"
        >
        <label for="photoInput" class="change-photo-button"> Modifica </label>
      </li>
      <li style="margin-left: 2rem">
        <strong>Nome Utente:</strong> {{ userName }}
        <input
          v-model="newName"
          placeholder="Modifica nome utente..."
          class="change-name-input"
        >
        <button class="change-name-button" @click="updateUserName">
          Modifica
        </button>
      </li>
      <li>
        <button class="return-button" @click="$router.push('/home')">
          Indietro
        </button>
      </li>
    </ul>
  </main>
</template>

<script setup>
import axios from "../services/axios.js";
import { userName, userPhoto } from "../services/auth.js";
import { isImageSet } from "../services/utilities.js";
import { ref } from "vue";

const serverUrl = __API_URL__;
const newName = ref("");
const newPhoto = userPhoto;

const updateUserName = async () => {
  try {
    const res = await axios.patch(`/users/me`, {
      name: newName.value,
    });
    sessionStorage.setItem("userName", res.data.name);
    userName.value = res.data.name;
  } catch (error) {}
};
const updateUserPhoto = async (file) => {
  try {
    const formData = new FormData();
    formData.append("photo", file);
    const res = await axios.patch(`/users/me/photo`, formData);
    sessionStorage.setItem("userPhoto", res.data.fileName);
    userPhoto.value = res.data.fileName;
  } catch (error) {}
};
</script>

<style scoped>
.dark-overlay {
  position: absolute;
  left: 0;
  bottom: 0;
  width: 37.05vw;
  height: calc(100vh - 3.7rem);
  background: black;
  opacity: 0.6;
  z-index: 1000;
}
.user-photo {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  margin-left: 0.5rem;
  border: 2px solid black;
  object-fit: cover;
  object-position: center;
}
.change-photo-button,
.change-name-button,
.return-button {
  width: fit-content;
  height: fit-content;
  margin-top: 3.5rem;
  margin-left: 1.5rem;
  padding: 0.2rem 0.4rem;
  background-color: whitesmoke;
  border: 2px solid black;
  border-color: black;
  color: black;
  cursor: pointer;
  font-size: 1rem;
}
.change-name-button:hover,
.return-button:hover,
.change-photo-button:hover {
  background-color: gainsboro;
}
.change-name-input {
  margin-left: 1rem;
  padding: 0.2rem;
  font-size: 1rem;
  background-color: whitesmoke;
  border: 1px solid black;
  outline: none;
}
.change-name-input::placeholder {
  font-size: 0.8rem;
  color: black;
  opacity: 0.7;
}
</style>
