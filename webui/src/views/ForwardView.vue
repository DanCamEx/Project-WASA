<template>
  <div class="dark-overlay" />
  <main>
    <h2 class="forward-title">Inoltra messaggio a...</h2>
    <ul class="conversations-container">
      <li
        v-for="conversation in conversations"
        :key="conversation.id"
        class="conversation"
        @click="getConversation(conversation.id)"
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
            ultimo messaggio: {{ conversation.latestMessage }}
          </div>
        </span>
      </li>
    </ul>
    <button class="return-button" @click="$router.back()">Cancella</button>
  </main>
</template>

<script setup>
import { onMounted, ref } from "vue";
import axios from "../services/axios.js";
import { useRoute } from "vue-router";
import { isImageSet } from "../services/utilities.js";

const route = useRoute();

const conversations = ref([]);
const getConversations = async () => {
	try {
		const res = await axios.get("/users/me/conversations");
		conversations.value = res.data;
	} catch (error) {}
};
const forwardMessage = async (conversationId) => {
	try {
		await axios.post(`/messages/${route.params.messageId}`, {
			conversationId: conversationId,
		});
		stopForward();
	} catch (error) {}
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
	position: absolute;
	right: 0;
	bottom: 0;
	width: 62.95vw;
	height: calc(100vh - 3.7rem);
	background: rgba(0, 0, 0, 0.6);
	z-index: 1000;
}
.return-button {
	margin-left: 18.5rem;
	margin-top: 1rem;
	font-size: 1rem;
	font-family: sans-serif;
	cursor: pointer;
	border-radius: 0;
	border-color: black;
	border-style: solid;
}
</style>

<style>
.forward-title {
	text-align: center;
	font-family: sans-serif;
	margin-top: 2rem;
	margin-right: 1rem;
}
.custom-list {
	margin-left: 4rem;
	margin-right: 4rem;
}
</style>
