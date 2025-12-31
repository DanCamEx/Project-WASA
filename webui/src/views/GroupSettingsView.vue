<template>
  <div class="dark-overlay" />
  <main style="padding: 1rem; font-family: sans-serif">
    <h1>Impostazioni Gruppo</h1>
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
              : '/images/default_group_photo.jpg'
          "
          alt="Foto Gruppo"
          width="150"
          height="150"
          class="group-photo"
        >
        <input
          id="photoInput"
          type="file"
          :hidden="true"
          accept="/*"
          @click="$event.target.value = null"
          @change="updateGroupPhoto($event.target.files[0])"
        >
        <label for="photoInput" class="change-photo-button"> Modifica </label>
      </li>
      <li style="margin-left: 2rem">
        <strong>Nome Gruppo:</strong> {{ currentName }}
        <input
          v-model="newName"
          placeholder="Modifica nome gruppo..."
          class="change-name-input"
        >
        <button class="change-name-button" @click="updateGroupName">
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
import { isImageSet } from "../services/utilities.js";
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";


const serverUrl = __API_URL__;
const currentName = ref("");
const newName = ref("");
const newPhoto = ref("");
const route = useRoute();
const conversationId = route.params.conversationId;

const updateGroupName = async () => {
  try {
    const res = await axios.patch(`/conversations/${conversationId}`, {
      name: newName.value,
    });
    currentName.value = res.data.name;
  } catch (error) {}
};
const updateGroupPhoto = async (file) => {
  try {
    const formData = new FormData();
    formData.append("photo", file);
    const res = await axios.patch(`/conversations/${conversationId}/photo`, formData);
    newPhoto.value = res.data.fileName;
  } catch (error) {}
};
const getGroupInfo = async () => {
  try {
    const res = await axios.get(`/conversations/${conversationId}`);
    currentName.value = res.data.name;
    newPhoto.value = res.data.photo;
  } catch (error) {}
};

onMounted(() => {
  getGroupInfo();
});
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
.group-photo {
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
  font-family: sans-serif;
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
  font-family: sans-serif;
  font-size: 0.8rem;
  color: black;
  opacity: 0.7;
}
</style>
