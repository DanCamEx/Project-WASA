<template>
  <UserSearchBar @search="searchUsers" />
  <ul>
    <li v-for="user in searchedUsers" :key="user.id">
      <span style="margin-left: 1rem; font-weight: bold">{{ user.name }}</span>
      <button @click="createConversation(user.id, 'direct')">
        Crea conversazione
      </button>
    </li>
  </ul>
  <ConversationsList
    :conversations="conversations"
    :latest-messages-content="latestMessagesContent"
    :get-conversation="getConversation"
    :get-latest-message="getLatestMessage"
  />
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import axios from "../services/axios.js";
import ConversationsList from "../components/ConversationsList.vue";
import UserSearchBar from "../components/UserSearchBar.vue";

const router = useRouter();

// Search Bar vars and methods
const searchedUsers = ref([]);
const allUsers = ref([]);
const searchUsers = (query) => {
  const search = query.toLowerCase();
  searchedUsers.value = allUsers.value.filter((user) =>
    user.name.toLowerCase().includes(search)
  );
  if (!search) {
    searchedUsers.value = [];
  }
};
const getAllUsers = async () => {
  try {
    const res = await axios.get("/users/search");
    allUsers.value = res.data;
  } catch (error) {
    console.error("Errore nel recupero degli utenti:", error);
  }
};

// Conversations List vars and methods
const conversations = ref([]);
const latestMessagesContent = ref({});
const getConversation = async (id) => {
  router.push("/conversations/" + id);
};
const createConversation = async (userId, category) => {
  try {
    const res = await axios.post("/users/me/conversations", {
      participantIds: [userId],
      category: category,
    });
    conversations.value.push(res.data);
  } catch (error) {
    console.error("Errore nella creazione della conversazione:", error);
  }
};
const getConversations = async () => {
  try {
    const res = await axios.get("/users/me/conversations");
    conversations.value = res.data;
    for (const conversation of conversations.value) {
      getLatestMessage(conversation.latestMessageId, conversation.id);
    }
  } catch (error) {
    console.error("Errore nel recupero delle conversazioni:", error);
  }
};
const getLatestMessage = async (id, conversationId) => {
  if (!id) {
    latestMessagesContent.value[conversationId] = "";
    return;
  }
  try {
    const res = await axios.get("/messages/" + id);
    latestMessagesContent.value[conversationId] = res.data.content;
  } catch (error) {
    console.error("Errore nel recupero dell'ultimo messaggio:", error);
  }
};

onMounted(async () => {
  await getConversations();
  await getAllUsers();
});
</script>

<style scoped>
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  width: auto;
  padding-right: 0.5rem;
  margin-top: 0.5rem;
  margin-left: 1rem;
  margin-right: -1rem;
  border: 0.2rem solid #000000;
  background-color: #f1f1f1;
  font-family: sans-serif;
  cursor: pointer;
}
button {
  font-size: 1rem;
  font-family: sans-serif;
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
  margin-left: 2rem;
  cursor: pointer;
  border-radius: 0;
  border-color: #000000;
  border-style: solid;
}
</style>
