<template>
  <main style="padding: 1rem">
    <div class="dark-overlay" />
    <h2>Nuovo Gruppo</h2>
    <ul class="group-creation-form">
      <strong class="group-photo-title">Foto Gruppo</strong>
      <li style="display: flex; flex-direction: row">
        <img
          :src="
            isImageSet(groupPhotoPreview)
              ? groupPhotoPreview
              : '/images/default_group_photo.jpg'
          "
          alt="Anteprima Foto Gruppo"
          width="150"
          height="150"
          class="group-photo-preview"
        >
        <input
          id="groupPhotoInput"
          type="file"
          accept="image/*"
          :hidden="true"
          @click="$event.target.value = null"
          @change="setTempGroupPhoto($event.target.files[0])"
        >
        <label
          class="choose-photo-button"
          style="margin-top: 3.5rem; margin-left: 1rem"
          for="groupPhotoInput"
        >Scegli Foto</label>
      </li>
      <li>
        Nome Gruppo:
        <input
          v-model="groupName"
          class="choose-name-input"
          placeholder="Inserisci il nome del gruppo..."
        >
      </li>
      <li>
        Membri:
        <input
          v-model="participantsNames"
          class="choose-name-input"
          style="margin-top: 1rem"
          placeholder="Inserisci i nomi dei membri..."
        >
        !!! I nomi devono essere separati da virgola e spazio (Es: Topolino,
        Paperino, Pippo)
      </li>
      <li>
        <button
          class="create-group-button"
          style="margin-top: 0.5rem"
          @click="createGroup()"
        >
          Crea Gruppo
        </button>
      </li>
      <li>
        <button
          class="return-button"
          style="margin-top: 3rem"
          @click="$router.push('/home')"
        >
          Indietro
        </button>
      </li>
    </ul>
  </main>
</template>

<script setup>
import axios from "../services/axios.js";
import { ref } from "vue";
import { isImageSet } from "../services/utilities.js";

const groupPhotoPreview = ref(null);
const groupPhotoFile = ref(null);
const groupName = ref("");
const participantsNames = ref("");
const participantsIds = ref([]);
const groupId = ref(0);

function setTempGroupPhoto(file) {
  if (!file) return;
  else {
    groupPhotoPreview.value = URL.createObjectURL(file);
    groupPhotoFile.value = file;
  }
}
const getParticipantsIds = async (names) => {
  participantsIds.value = [];
  for (const name of names.split(", ")) {
    try {
      const res = await axios.get("/users/by-name/" + name);
      participantsIds.value.push(res.data.id);
    } catch (error) {}
  }
};
const createGroup = async () => {
  await getParticipantsIds(participantsNames.value);
  if (groupName.value.trim() === "") {
    groupName.value = "Nuovo Gruppo";
  }
  try {
    const res = await axios.post("/users/me/conversations", {
      name: groupName.value,
      participantsIds: participantsIds.value,
      category: "group",
    });
    groupId.value = res.data.id;
  } catch (error) {
    return;
  }
  try {
    const formData = new FormData();
    formData.append("photo", groupPhotoFile.value);
    const res = await axios.patch(
      "/conversations/" + groupId.value + "/photo",
      formData
    );
  } catch (error) {
    return;
  }
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
.group-creation-form {
  display: flex;
  flex-direction: column;
  list-style-type: none;
}
.group-photo-title {
  display: flex;
  flex-direction: column;
  margin-left: 1.9rem;
  margin-bottom: 0.3rem;
}
.group-photo-preview {
  border-radius: 50%;
  border: 3px solid black;
  object-fit: cover;
  object-position: center;
}
.create-group-button,
.choose-photo-button,
.return-button {
  width: fit-content;
  height: fit-content;
  padding: 0.2rem 0.4rem;
  background-color: whitesmoke;
  border: 2px solid black;
  border-color: black;
  color: black;
  cursor: pointer;
  font-size: 1rem;
}
.create-group-button:hover,
.return-button:hover,
.choose-photo-button:hover {
  background-color: gainsboro;
}
.choose-name-input {
  margin-left: 1rem;
  padding: 0.2rem;
  font-size: 1rem;
  background-color: whitesmoke;
  border: 1px solid black;
  outline: none;
}
.choose-name-input::placeholder {
  font-size: 0.8rem;
  color: black;
  opacity: 0.7;
}
</style>
