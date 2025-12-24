<template>
  <input
    type="text"
    placeholder="Scrivi un messaggio..."
    v-model="sendQuery"
    class="send-bar"
    @keyup.enter="sendMessage(sendQuery.value)"
  />
</template>

<script setup>
import { ref } from "vue";
const sendQuery = ref("");
const emit = defineEmits(["send"]);
const sendMessage = () => {
  emit("send", sendQuery.value);
  sendQuery.value = "";
};
</script>

<style scoped>
.send-bar {
  width: 99%;
  margin-left: 0.5rem;
  box-sizing: border-box;
  height: 3rem;
  border: 0.1rem solid #000000;
  outline: none;
  font-size: 1rem;
  background-color: #f1f1f1;
  font-family: sans-serif;
  font-size: 1.2rem;
  color: #000000;
}
input::placeholder {
  font-family: sans-serif;
  font-size: 1.2rem;
  color: #000000;
  font-style: italic;
}
</style>
