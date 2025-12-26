<template>
	<div class="messages">
		<div
			v-for="message in messages"
			:key="message.id"
			:class="
				message.senderId == authId ? 'message-right' : 'message-left'
			"
			:style="
				message.senderId == authId
					? 'align-self: flex-end;'
					: 'align-self: flex-start;'
			"
			@click="
				showMessageOptions =
					showMessageOptions == message.id ? null : message.id;
				showEmotes = false;
			"
		>
			{{ message.content }}
			<span v-if="message.senderId == authId">
				({{ message.status }})</span
			>
			<span
				class="options"
				v-if="showMessageOptions == message.id"
				:style="
					message.senderId == authId
						? 'right:0;left:auto;'
						: 'left:0;right:auto;'
				"
			>
				<div v-if="showEmotes">
					<button
						class="options-button"
						@click="commentMessage(message.id, '👍')"
					>
						👍
					</button>
					<button
						class="options-button"
						@click="commentMessage(message.id, '😂')"
					>
						😂
					</button>
					<button
						class="options-button"
						@click="commentMessage(message.id, '💀')"
					>
						💀
					</button>
					<button
						class="options-button"
						@click="commentMessage(message.id, '😭')"
					>
						😭
					</button>
					<button
						class="options-button"
						@click="commentMessage(message.id, '🔥')"
					>
						🔥
					</button>
				</div>
				<button
					class="options-button"
					@click.stop="showEmotes = !showEmotes"
				>
					Commenta
				</button>
				<button
					class="options-button"
					@click="forwardMessage(message.id)"
				>
					Inoltra
				</button>
				<button
					class="options-button"
					v-if="
						(comments[message.id] || []).some(
							(c) => c.reacterId == authId
						)
					"
					@click="unCommentMessage(message.id)"
				>
					Elimina commento
				</button>
			</span>
			<span
				v-for="(comment, idx) in comments[message.id] || []"
				:key="message.id + '-' + comment.emoji"
				class="emoji"
				:style="`left: ${idx * 1.1}rem;`"
			>
				{{ comment.emoji }}
			</span>
			<span class="forwarded" v-if="message.isForwarded == true">
				Inoltrato
			</span>
		</div>
	</div>
	<SendMessageBar @send="sendMessage" />
</template>

<script setup>
import { onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import axios from "../services/axios.js";
import SendMessageBar from "../components/SendMessageBar.vue";
import router from "../router/index.js";

const route = useRoute();
const authId = localStorage.getItem("authId");

const showEmotes = ref(false);
const showMessageOptions = ref(null);

const getConversation = async (id) => {
	try {
		const res = await axios.get("/conversations/" + id + "/messages");
		if (res.data == null) {
			messages.value = [];
			return;
		}
		messages.value = res.data.reverse();
		for (const message of messages.value) {
			getComments(message.id);
		}
	} catch (error) {}
};

const messages = ref([]);
const sendMessage = async (content) => {
	try {
		const res = await axios.post(
			"/conversations/" + route.params.id + "/messages",
			{
				content: content,
			}
		);
		messages.value.unshift(res.data);
	} catch (error) {}
};

const comments = ref({});
const commentMessage = async (messageId, emoji) => {
	try {
		const res = await axios.post("/messages/" + messageId + "/comments", {
			emoji: emoji,
		});
		getComments(messageId);
	} catch (error) {}
};
const unCommentMessage = async (messageId) => {
	try {
		const res = await axios.delete("/messages/" + messageId + "/comments");
		getComments(messageId);
	} catch (error) {}
};
const getComments = async (messageId) => {
	try {
		const res = await axios.get("/messages/" + messageId + "/comments");
		comments.value[messageId] = res.data;
	} catch (error) {}
};

const forwardMessage = async (id) => {
	router.push("/messages/" + id + "/forward");
};

onMounted(() => {
	getConversation(route.params.id);
});

watch(
	() => route.params.id,
	(newId) => {
		getConversation(newId);
	}
);
</script>

<style scoped>
.messages {
	height: calc(100% - 4.1rem);
	display: flex;
	flex-direction: column-reverse;;
	overflow-y: auto;
	padding: 0.5rem;
}
.message-right,
.message-left {
	position: relative;
	max-width: 60%;
	padding: 0.5rem;
	margin: 1rem 0rem;
	border-radius: 0.1rem;
	font-family: sans-serif;
	background-color: #f1f1f1;
	word-break: break-word;
	cursor: pointer;
}
.message-left:hover,
.message-right:hover {
	background-color: #e0e0e0;
}
.options {
	position: absolute;
	white-space: nowrap;
	bottom: 100%;
	z-index: 1;
}
.options-button {
	font-size: 1rem;
	font-family: sans-serif;
	margin-right: 0.5rem;
	margin-top: 0.5rem;
	cursor: pointer;
	border-radius: 0;
	border-color: #000000;
	border-style: solid;
}
.emoji {
	position: absolute;
	bottom: -20%;
	font-size: 0.9rem;
}
.forwarded {
	position: absolute;
	bottom: 100%;
	right: 0;
	font-size: 0.9rem;
	font-family: sans-serif;
	font-weight: lighter;
	background-color: #f1f1f1;
	padding: 0.1rem 0.4rem;
}
</style>
