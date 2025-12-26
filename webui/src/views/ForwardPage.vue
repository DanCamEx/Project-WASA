<template>
	<div class="dark-overlay"></div>
	<ConversationsList
		:title="'A chi vuoi inoltrare il messaggio?'"
		titleClass="conversations-title"
		listClass="customList"
		:forwardMessage="forwardMessage"
		:conversations="conversations"
		:latestMessagesContent="latestMessagesContent"
		:getConversation="getConversation"
		:getLatestMessage="getLatestMessage"
	/>
	<button @click="stopForward" :style="{ cursor: 'pointer' }">
		Cancella
	</button>
</template>

<script setup>
import { onMounted, ref } from "vue";
import ConversationsList from "../components/ConversationsList.vue";
import axios from "../services/axios.js";
import { useRoute } from "vue-router";

const route = useRoute();

const conversations = ref([]);
const getConversations = async () => {
	try {
		const res = await axios.get("/users/me/conversations");
		conversations.value = res.data;
	} catch (error) {
		console.error("Errore nel recupero delle conversazioni:", error);
	}
};
const forwardMessage = async (conversationId) => {
	try {
		await axios.post(`/messages/${route.params.id}`, {
			conversationId: conversationId,
		});
		stopForward();
	} catch (error) {
		console.error("Errore nell'inoltro del messaggio:", error);
	}
};
const stopForward = () => {
	history.back();
};

onMounted(async () => {
	await getConversations();
});

// This is blank because the logic is handled in HomePage.vue
const latestMessagesContent = ref({});
const getLatestMessage = () => {};
const getConversation = () => {};

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
button {
	margin-left: 19rem;
	margin-top: 1rem;
	font-size: 1rem;
	font-family: sans-serif;
	cursor: pointer;
	border-radius: 0;
	border-color: #000000;
	border-style: solid;
}
</style>

<style>
.coneversations-title {
	text-align: center;
	font-family: sans-serif;
	margin-top: 1rem;
}
.customList {
	margin-left: 4rem;
	margin-right: 4rem;
}
</style>
