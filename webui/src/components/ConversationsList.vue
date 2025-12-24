<template>
  <h2 class="conversations-title">{{ title }}</h2>
  <ul :class="customClass">
    <li
      v-for="conversation in conversations"
      :key="conversation.id"
      @click="
        getConversation(conversation.id);
        getLatestMessage(conversation.latestMessageId, conversation.id);
        forwardMessage(conversation.id);
      "
    >
      <div class="conversation-name">{{ conversation.name }}</div>
      <div class="latest-message" v-if="latestMessagesContent[conversation.id]">
        ultimo messaggio:{{ latestMessagesContent[conversation.id] }}
      </div>
    </li>
  </ul>
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    default: "Le tue conversazioni",
  },
  customClass: {
    type: String,
    default: "",
  },
  conversations: Array,
  latestMessagesContent: Object,
  getConversation: Function,
  getLatestMessage: Function,
  forwardMessage: Function,
});
</script>

<style scoped>
ul {
  list-style-type: none;
  padding: 0;
}
li {
  padding: 0.5rem;
  margin-left: 1rem;
  margin-right: 1rem;
  margin-bottom: 0.5rem;
  border: 0.2rem solid #000000;
  background-color: #f1f1f1;

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

<style>
.conversations-title {
  margin-left: 1rem;
  font-size: 1.5rem;
  font-family: sans-serif;
}
</style>
