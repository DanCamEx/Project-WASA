<template>
  <main>
    <input
      v-model="usersToSearch"
      type="text"
      placeholder="Cerca un utente..."
      class="search-bar"
      @input="searchUsers(usersToSearch)"
    >
    <ul class="users-container">
      <li v-for="user in usersFound" :key="user.id" class="user-tab">
        <img
          width="50"
          height="50"
          :src="
            isImageSet(user.photo)
              ? serverUrl + '/images/' + user.photo
              : '/images/default_user_photo.jpg'
          "
          alt="Foto utente"
          style="border-radius: 50%; object-fit: cover; object-position: center"
        >
        <div style="margin-left: 0.3rem; font-weight: bold">
          {{ user.name }}
        </div>
        <button
          class="create-direct-button"
          @click="createDirectConversation(user.id, 'direct')"
        >
          Inizia una conversazione
        </button>
      </li>
    </ul>
    <div class="conversations-header">
      <h2 class="conversations-title">Le tue conversazioni</h2>
      <button class="create-group-button" @click="router.push('/new-group')">
        Crea gruppo
      </button>
    </div>
    <ul class="conversations-container">
      <li
        v-for="conversation in conversations"
        :key="conversation.id"
        class="conversation"
        @click="getConversation(conversation.id)"
      >
        <span>
          <img
            :src="
              isImageSet(conversation.photo)
                ? serverUrl + '/images/' + conversation.photo
                : getDefaultConversationImage(conversation.category)
            "
            class="conversation-photo"
            alt="Foto conversazione"
          >
        </span>
        <span class="conversation-info">
          <span class="conversation-name">{{ conversation.name }}</span>
          <div
            v-if="conversation.latestMessage.content != ''"
            class="latest-message"
          >
            <template
              v-if="
                conversation.latestMessage.content.endsWith(
                  '.jpg' || '.png' || '.jpeg' || '.gif'
                )
              "
            >
              ultimo messaggio:
              <img class="icon" src="/icons/image.svg">
              <br>
            </template>
            <template v-else>ultimo messaggio: {{ conversation.latestMessage.content }}<br></template>
            inviato il:
            {{ formatTimestamp(conversation.latestMessage.timestamp) }}
          </div>
          <span v-if="conversation.category == 'group'" class="group-tag">gruppo</span>
          <span v-else class="direct-tag">diretta</span>
        </span>
      </li>
    </ul>
  </main>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import axios from "../services/axios.js";
import {
  isImageSet,
  getDefaultConversationImage,
  formatTimestamp,
  startPolling
} from "../services/utilities.js";

const router = useRouter();
const serverUrl = __API_URL__;

const allUsers = ref([]);
const getAllUsers = async () => {
  try {
    const res = await axios.get("/users/search");
    allUsers.value = res.data;
  } catch (error) {}
};

const usersToSearch = ref("");
const usersFound = ref([]);
const searchUsers = (usersToSearch) => {
  const searching = usersToSearch.toLowerCase();
  usersFound.value = allUsers.value.filter((user) =>
    user.name.toLowerCase().includes(searching)
  );
  if (!searching) {
    usersFound.value = [];
  }
};

const getConversation = async (conversationId) => {
  try {
    router.push("/conversations/" + conversationId);
  } catch (error) {}
};

const conversations = ref([]);
const createDirectConversation = async (userId, category) => {
  try {
    const res = await axios.post("/users/me/conversations", {
      name: "",
      participantsIds: [userId],
      category: category,
    });
    conversations.value.unshift(res.data);
  } catch (error) {}
};
const getConversations = async () => {
  try {
    const res = await axios.get("/users/me/conversations");
    if (res.data == null) {
      conversations.value = [];
      return;
    }
    conversations.value = res.data;
    conversations.value.sort(
      (a, b) =>
        new Date(b.latestMessage.timestamp) -
        new Date(a.latestMessage.timestamp)
    );
  } catch (error) {}
};

onMounted(async () => {
  await getConversations();
  await getAllUsers();
  startPolling(getConversations);
});
</script>

<style scoped>
.search-bar {
  width: 95%;
  height: 2rem;
  margin-left: 0.8rem;
  margin-top: 1rem;
  border: 0.1rem solid black;
  font-family: sans-serif;
  font-size: 1rem;
  font-style: italic;
  background-color: whitesmoke;
  color: black;
  outline: none;
}
.search-bar::placeholder {
  color: black;
}
.users-container {
  list-style-type: none;
  margin: 0;
  padding: 0;
  overflow-y: auto;
}
.user-tab {
  display: flex;
  align-items: center;
  width: 50%;
  margin-left: 3rem;
  margin-top: 0.5rem;
  padding: 0.5rem;
  border: 0.2rem solid black;
  background-color: whitesmoke;
  font-family: sans-serif;
  cursor: pointer;
}
.create-direct-button {
  width: fit-content;
  margin-left: 2rem;
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
  border-radius: 0;
  border-color: black;
  border-style: solid;
  font-size: 1rem;
  font-family: sans-serif;
  cursor: pointer;
}
.conversations-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  margin-left: 1rem;
  margin-top: 1rem;
}
.create-group-button {
  min-width: 120px;
  margin-left: 2rem;
  margin-right: 1rem;
  margin-bottom: 0.2rem;
  border-radius: 0;
  border-color: black;
  border-style: solid;
  font-size: 1rem;
  font-family: sans-serif;
  cursor: pointer;
}
.conversations-title {
  width: fit-content;
  margin-top: 0;
  margin-bottom: 0.2rem;
  font-size: 1.5rem;
  font-family: sans-serif;
}
.conversations-container {
  list-style-type: none;
  padding: 0;
  margin: 0;
}
.conversation {
  padding: 0.5rem;
  margin-left: 1rem;
  margin-right: 1rem;
  margin-top: 0.5rem;
  border: 0.2rem solid black;
  background-color: whitesmoke;
  font-family: sans-serif;
  cursor: pointer;
  display: flex;
  align-items: center;
}
.conversation:hover {
  background-color: gainsboro;
}
.conversation-info {
  display: flex;
  flex-direction: column;
  margin-left: 1rem;
  overflow: hidden;
}
.direct-tag,
.group-tag {
  font-size: 0.8rem;
  font-weight: bold;
  color: white;
  background-color: black;
  padding: 0.1rem 0.3rem;
  width: fit-content;
}
.conversation-name {
  font-size: 1.2rem;
  font-weight: bold;
}
.conversation-photo {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  border: 3px solid black;
  object-fit: cover;
  object-position: center;
}
.latest-message {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: black;
  font-size: 0.9rem;
}
.icon {
  width: 1.2em;
  height: 1.2em;
  vertical-align: middle;
  stroke: black;
  fill: none;
  stroke-width: 2;
}
</style>
