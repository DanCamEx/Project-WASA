<template>
  <div class="dark-overlay"></div>
  <ConversationsList
    :title="'&nbsp &nbsp &nbsp &nbsp &nbsp &nbsp &nbsp A chi vuoi inoltrare il messaggio?'"
    customClass="customList"
    :forwardMessage="forwardMessage"
    :conversations="conversations"
    :latest-messages-content="{}"
  />
  <button @click="stopForward" :style="{ cursor: 'pointer' }">Cancella</button>
</template>

<script setup>
import { onMounted, ref } from "vue";
import ConversationsList from "../components/ConversationsList.vue";
import axios from "../services/axios.js";

const conversations = ref([]);
const getConversations = async () => {
  try {
    const res = await axios.get("/users/me/conversations");
    conversations.value = res.data;
    console.log(res.data);
  } catch (error) {
    console.error("Errore nel recupero delle conversazioni:", error);
  }
};

const stopForward = () => {
  history.back();
};

onMounted(async () => {
  await getConversations();
});
</script>

<style scoped>
.dark-overlay {
  position: fixed;
  top: 0;
  right: 0;
  width: 66.66vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.6);
  z-index: 1000;
}
button {
  margin-left: 16rem;
  margin-top: 1rem;
  font-size: 1rem;
  font-family: sans-serif;
}
</style>

<style>
    .customList {
        margin-left: 4rem;
        margin-right: 4rem;
    }
</style>
