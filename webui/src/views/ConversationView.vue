<template>
  <header class="conversation-header">
    <img
      :src="
        isImageSet(conversationInfo.photo)
          ? serverUrl + '/images/' + conversationInfo.photo
          : getDefaultConversationImage(conversationInfo.category)
      "
      class="conversation-photo-header"
      alt="Foto conversazione"
    />
    <h1 class="conversation-title">{{ conversationInfo.name }}</h1>
    <div
      v-if="conversationInfo.category == 'group'"
      class="conversation-actions"
    >
      <input
        v-model="newMemberName"
        type="text"
        placeholder="Aggiungi membro al gruppo..."
        class="add-to-group-input"
      />
      <button class="add-to-group-button" @click="addToGroup">Aggiungi</button>
      <button class="leave-group-button" @click="leaveGroup">
        Esci dal gruppo
      </button>
      <img
        class="settings-icon"
        src="/icons/gear.svg"
        @click="
          router.push('/conversations/' + conversationInfo.id + '/settings')
        "
      />
    </div>
  </header>
  <ul class="messages">
    <li
      v-for="message in messages"
      :key="message.id"
      :class="message.senderId == authId ? 'message-right' : 'message-left'"
      :style="
        message.senderId == authId
          ? 'align-self: flex-end;'
          : 'align-self: flex-start;'
      "
    >
      <span v-if="message.senderId == authId" class="message-timestamp">
        {{ formatTimestamp(message.timestamp) }}
      </span>
      <span
        class="message-content"
        @click="
          showMessageOptions =
            showMessageOptions == message.id ? null : message.id;
          showEmotes = false;
        "
        ><img
          v-if="
            message.content.endsWith('.jpg') ||
            message.content.endsWith('.png') ||
            message.content.endsWith('.jpeg') ||
            message.content.endsWith('.gif')
          "
          :src="serverUrl + '/images/' + message.content"
          alt="Immagine inviata"
          style="max-width: 300px; max-height: 300px"
        />
        <span v-else>
          {{ message.content }}
        </span>
        <template v-if="message.senderId == authId && message.status == 'sent'">
          <img class="check-icon" src="/icons/check2.svg" />
        </template>
        <template
          v-else-if="message.senderId == authId && message.status == 'read'"
        >
          <img class="check-icon" src="/icons/check2-all.svg" />
        </template>
      </span>
      <span v-if="message.senderId != authId" class="sender-name"
        >{{ sendersNames[message.senderId] || "" }} -
      </span>
      <span v-if="message.senderId != authId" class="message-timestamp">
        {{ formatTimestamp(message.timestamp) }}
      </span>
      <span
        v-if="showMessageOptions == message.id"
        class="options"
        :style="
          message.senderId == authId
            ? 'right:0;left:auto;'
            : 'left:0;right:auto;'
        "
      >
        <template v-if="showEmotes">
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
        </template>
        <button class="options-button" @click.stop="showEmotes = !showEmotes">
          Commenta
        </button>
        <button class="options-button" @click="forwardMessage(message.id)">
          Inoltra
        </button>
        <button
          v-if="(comments[message.id] || []).some((c) => c.reacterId == authId)"
          class="options-button-delete"
          @click="unCommentMessage(message.id)"
        >
          Elimina commento
        </button>
        <button
          v-if="message.senderId == authId"
          class="options-button-delete"
          @click="deleteMessage(message.id)"
        >
          Elimina messaggio
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
      <span v-if="message.isForwarded == true" class="forwarded">
        Inoltrato
      </span>
    </li>
  </ul>
  <footer class="send-footer">
    <div class="send-bar-wrapper">
      <input
        v-model="messageToSend"
        type="text"
        placeholder="Scrivi un messaggio..."
        class="send-bar"
        @keyup.enter="sendMessage"
      />
      <label for="imageInput" class="send-photo-button">Invia Immagine</label>
    </div>
    <input
      id="imageInput"
      type="file"
      hidden
      accept="image/*"
      @change="sendImage($event.target.files[0])"
    />
  </footer>
</template>

<script setup>
import { onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import axios from "../services/axios.js";
import router from "../router/index.js";
import { userAuthId } from "../services/auth.js";
import {
  formatTimestamp,
  isImageSet,
  getDefaultConversationImage,
} from "../services/utilities.js";

const serverUrl = __API_URL__;
const route = useRoute();
const messageToSend = ref("");
const authId = userAuthId;

const conversationInfo = ref({});
const showEmotes = ref(false);
const showMessageOptions = ref(null);

const getConversation = async (conversationId) => {
  try {
    const res = await axios.get("/conversations/" + conversationId);
    conversationInfo.value = res.data;
  } catch (error) {}
};
const getMessages = async (conversationId) => {
  try {
    const res = await axios.get(
      "/conversations/" + conversationId + "/messages"
    );
    if (res.data == null) {
      messages.value = [];
      return;
    }
    messages.value = res.data.reverse();
    for (const message of messages.value) {
      getComments(message.id);
    }
    getSendersNames();
    console.log(messages.value);
  } catch (error) {}
};

const messages = ref([]);
const sendMessage = async () => {
  try {
    if (messageToSend.value == "") {
      return;
    }
    const res = await axios.post(
      "/conversations/" + route.params.conversationId + "/messages",
      {
        content: messageToSend.value,
      }
    );
    messages.value.unshift(res.data);
    messageToSend.value = "";
  } catch (error) {}
};

const sendImage = async (imageFile) => {
  try {
    const formData = new FormData();
    formData.append("image", imageFile);
    const res = await axios.post(
      "/conversations/" + route.params.conversationId + "/messages/image",
      formData
    );
    messages.value.unshift(res.data);
  } catch (error) {}
};
const deleteMessage = async (messageId) => {
  try {
    await axios.delete(
      "/conversations/" + route.params.conversationId + "/messages/" + messageId
    );
    messages.value = messages.value.filter((m) => m.id != messageId);
  } catch (error) {}
};
const sendersNames = ref({});
const senderIds = ref([]);
const getSendersNames = async () => {
  senderIds.value = [
    ...new Set(messages.value.map((message) => message.senderId)),
  ];
  for (const senderId of senderIds.value) {
    try {
      const res = await axios.get("/users/by-id/" + senderId);
      sendersNames.value[senderId] = res.data.name;
    } catch (error) {}
  }
  for (const message of messages.value) {
    message.senderName = sendersNames.value[message.senderId];
  }
};

const comments = ref({});
const commentMessage = async (messageId, emoji) => {
  try {
    await axios.post("/messages/" + messageId + "/comments", {
      emoji: emoji,
    });
    getComments(messageId);
  } catch (error) {}
};

const unCommentMessage = async (messageId) => {
  try {
    await axios.delete("/messages/" + messageId + "/comments");
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

const newMemberName = ref("");
const addToGroup = async () => {
  try {
    await axios.post(
      "/conversations/" + route.params.conversationId + "/participants",
      {
        userName: newMemberName.value,
      }
    );
    newMemberName.value = "";
  } catch (error) {}
};
const leaveGroup = async () => {
  try {
    await axios.delete(
      "/conversations/" + route.params.conversationId + "/participants"
    );
    router.push("/home");
  } catch (error) {}
};
onMounted(() => {
  getMessages(route.params.conversationId);
  getConversation(route.params.conversationId);
});

watch(
  () => route.params.conversationId,
  (newConversationId) => {
    getMessages(newConversationId);
    getConversation(newConversationId);
  }
);
</script>

<style scoped>
.conversation-header {
  display: flex;
  align-items: center;
  border-bottom: solid black;
  padding-bottom: 0.3rem;
  background-color: whitesmoke;
  position: sticky;
  top: 0;
  z-index: 1;
  overflow-y: hidden;
}
.conversation-title {
  margin-left: 1.5rem;
  margin-top: 1.6rem;
  font-family: sans-serif;
}
.conversation-actions {
  display: flex;
  align-items: center;
  margin-left: auto;
}
.conversation-photo-header {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 50%;
  border: 0.2rem solid black;
  margin-top: 0.2rem;
  margin-left: 1rem;
}
.conversation-actions {
  justify-content: flex-end;
  margin-right: 1rem;
}
.add-to-group-input {
  margin-top: 2rem;
  height: 1.5rem;
  width: 15rem;
  border: 0.1rem solid black;
  font-family: sans-serif;
  font-size: 0.9rem;
  background-color: whitesmoke;
  color: black;
  outline: none;
}
.add-to-group-input::placeholder {
  color: black;
}
.add-to-group-button {
  width: fit-content;
  margin-left: 0.2rem;
  margin-top: 2rem;
  height: 2rem;
  border-color: black;
  border-style: solid;
  font-size: 1rem;
  font-family: sans-serif;
  cursor: pointer;
}
.leave-group-button {
  width: fit-content;
  margin-left: 1rem;
  margin-top: 2rem;
  height: 2.2rem;
  border-color: darkred;
  border-style: solid;
  color: darkred;
  font-size: 1rem;
  font-family: sans-serif;
  cursor: pointer;
}
.messages {
  list-style-type: none;
  height: calc(100% - 10.5rem);
  display: flex;
  flex-direction: column-reverse;
  overflow-y: auto;
  padding: 0.5rem;
  margin: 0;
}
.message-right,
.message-left {
  position: relative;
  max-width: 60%;
  margin: 1.2rem 0rem 1rem 0rem;
  font-family: sans-serif;
  word-break: break-word;
  cursor: pointer;
  background: none;
  padding: 0;
}
.message-row {
  display: flex;
  align-items: center;
  gap: 0.3rem;
}
.message-timestamp {
  margin-right: 0.5rem;
  font-size: 0.8rem;
  cursor: default;
}
.sender-name {
  margin-left: 0.5rem;
  font-weight: bold;
  font-size: 0.95em;
  cursor: default;
  white-space: nowrap;
  background: none;
}
.message-content {
  background-color: whitesmoke;
  border-radius: 0.1rem;
  padding: 0.5rem;
  display: inline-block;
  min-width: 2.5rem;
  position: relative;
}
.message-content:hover {
  background-color: gainsboro;
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
  cursor: pointer;
  border-radius: 0;
  border-color: black;
  border-style: solid;
}
.options-button-delete {
  font-size: 1rem;
  font-family: sans-serif;
  margin-right: 0.5rem;
  cursor: pointer;
  border-radius: 0;
  border-color: darkred;
  border-style: solid;
  color: darkred;
}
.emoji {
  position: absolute;
  bottom: -0.5rem;
  font-size: 0.9rem;
}
.forwarded {
  position: absolute;
  bottom: 100%;
  right: 0;
  font-size: 0.9rem;
  font-family: sans-serif;
  font-weight: lighter;
  background-color: whitesmoke;
  padding: 0.1rem 0.4rem;
}
.send-footer {
  position: relative;
  padding: 0.2rem 0.2rem 0 0.2rem;
}
.send-bar-wrapper {
  position: relative;
}
.send-bar {
  width: 100%;
  box-sizing: border-box;
  height: 3rem;
  border: 0.1rem solid black;
  outline: none;
  font-size: 1rem;
  background-color: whitesmoke;
  font-family: sans-serif;
  font-size: 1.2rem;
  color: black;
  padding-right: 7rem;
}
.send-photo-button {
  position: absolute;
  right: 0.5rem;
  bottom: 0.3rem;
  z-index: 1;
}
.send-bar {
  width: 100%;
  box-sizing: border-box;
  height: 3rem;
  border: 0.1rem solid black;
  outline: none;
  font-size: 1rem;
  background-color: whitesmoke;
  font-family: sans-serif;
  font-size: 1.2rem;
  color: black;
}
.send-bar::placeholder {
  font-family: sans-serif;
  font-size: 1.2rem;
  color: black;
  font-style: italic;
}
.send-photo-button {
  width: fit-content;
  height: fit-content;
  margin-bottom: 0.3rem;
  padding: 0.2rem 0.4rem;
  background-color: whitesmoke;
  border: 0.2rem solid black;
  border-color: black;
  color: black;
  cursor: pointer;
  font-family: sans-serif;
  font-size: 1rem;
}
.send-photo-button:hover {
  background-color: gainsboro;
}
.check-icon {
  width: 0.9rem;
  height: 0.9rem;
  vertical-align: sub;
  stroke: black;
  fill: none;
  stroke-width: 2;
}
.settings-icon {
  width: 2rem;
  height: 2rem;
  margin-bottom: 3rem;
  stroke: black;
  fill: none;
  stroke-width: 2;
  cursor: pointer;
}
.settings-icon:hover {
  stroke: gray;
  background-color: darkgrey;
  border-radius: 50%;
}
</style>
