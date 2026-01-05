<template>
  <main>
    <div class="dark-overlay" />
    <h2 class="forward-title">Inoltra messaggio a...</h2>
    <ul class="conversations-container">
      <li
        v-for="conversation in conversations"
        :key="conversation.id"
        class="conversation"
        @click="
          getConversation(conversation.id), forwardMessage(conversation.id)
        "
      >
        <span class="conversation-photo">
          <img
            :src="
              isImageSet(conversation.image)
                ? serverUrl + '/images/' + conversation.image
                : getDefaultConversationImage(conversation.category)
            "
            alt="Foto conversazione"
            width="70"
            height="70"
            style="
              border-radius: 50%;
              object-fit: cover;
              object-position: center;
            "
          >
        </span>
        <span class="conversation-info">
          <span class="conversation-name">{{ conversation.name }}</span>
          <div v-if="conversation.latestMessage != ''" class="latest-message">
            ultimo messaggio: {{ conversation.latestMessage.content }}<br>
            inviato il
            {{ formatTimestamp(conversation.latestMessage.timestamp) }}
          </div>
          <span v-if="conversation.category == 'group'" class="group-tag">gruppo</span>
          <span v-else class="direct-tag">diretta</span>
        </span>
      </li>
    </ul>
    <button class="return-button" @click="$router.back()">Cancella</button>
  </main>
</template>

<script setup>
import { onMounted, ref } from "vue";
import axios from "../services/axios.js";
import { useRouter, useRoute } from "vue-router";
import {
  isImageSet,
  getDefaultConversationImage,
  formatTimestamp,
} from "../services/utilities.js";

const router = useRouter();
const route = useRoute();

const conversations = ref([]);
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
const getConversation = async (conversationId) => {
  try {
    router.push("/conversations/" + conversationId);
  } catch (error) {}
};
const forwardMessage = async (conversationId) => {
  try {
    await axios.post(`/messages/${route.params.messageId}/forward`, {
      conversationId: conversationId,
    });
  } catch (error) {}
};

onMounted(async () => {
  await getConversations();
});
</script>

<style scoped>
main {
  height: 100%;
  background-color: gainsboro;
  overflow-y: auto;
}
.dark-overlay {
  position: absolute;
  right: 0;
  bottom: 0;
  width: 62.95vw;
  height: calc(100vh - 3.7rem);
  background: rgba(0, 0, 0, 0.6);
  z-index: 1;
}
.return-button {
  margin-left: 18.5rem;
  margin-top: 1rem;
  font-size: 1rem;
  cursor: pointer;
  border-radius: 0;
  border-color: black;
  border-style: solid;
}
.conversations-container {
  list-style-type: none;
  padding: 0;
  margin: 0;
}
.conversation {
  display: flex;
  align-items: center;
  margin-left: 1rem;
  margin-right: 1rem;
  margin-top: 0.5rem;
  padding: 0.5rem;
  border: 0.2rem solid black;
  background-color: rgb(0, 150, 250);
  cursor: pointer;
}
.conversation:hover {
  background-color: rgb(0, 100, 200);
}
.conversation-info {
  display: flex;
  flex-direction: column;
  margin-left: 1rem;
  overflow: hidden;
}
.direct-tag,
.group-tag {
  width: fit-content;
  padding: 0.1rem 0.3rem;
  border: none;
  background-color: black;
  color: white;
  font-size: 0.8rem;
  font-weight: bold;
}
.conversation-name {
  font-size: 1.2rem;
  font-weight: bold;
}
.conversation-photo {
  width: 70px;
  height: 70px;
  border: 3px solid black;
  border-radius: 50%;
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
  stroke: black;
  stroke-width: 2;
  fill: none;
  vertical-align: middle;
}
</style>
