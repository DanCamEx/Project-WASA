<template>
	<h2 :class="['conversations-title', titleClass]">{{ title }}</h2>
	<ul :class="['conversations-container', listClass]">
		<li
			v-for="conversation in conversations"
			:key="conversation.id"
			@click="
				getConversation(conversation.id);
				getLatestMessage(conversation.latestMessageId, conversation.id);
				forwardMessage(conversation.id);
			"
			class="conversation"
		>
			<div class="conversation-name">{{ conversation.name }}</div>
			<div
				class="latest-message"
				v-if="latestMessagesContent[conversation.id]"
			>
				ultimo messaggio:{{ latestMessagesContent[conversation.id] }}
			</div>
		</li>
	</ul>
</template>

<script setup>
const props = defineProps({
	title: String,
	titleClass: String,
	listClass: String,
	conversations: Array,
	latestMessagesContent: Object,
	getConversation: Function,
	getLatestMessage: Function,
	forwardMessage: Function,
});
</script>

<style scoped>
.conversations-title {
	margin-left: 1rem;
	font-size: 1.5rem;
	font-family: sans-serif;
}
.conversations-container {
	list-style-type: none;
	padding: 0;
}
.conversation {
	padding: 0.5rem;
	margin-left: 1rem;
	margin-right: 1rem;
	margin-bottom: 0.5rem;
	border: 0.2rem solid #000000;
	background-color: #f1f1f1;
	font-family: sans-serif;
	cursor: pointer;
}
.conversation-name {
	font-weight: bold;
}
.latest-message {
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	color: rgb(0, 0, 0);
}
</style>
