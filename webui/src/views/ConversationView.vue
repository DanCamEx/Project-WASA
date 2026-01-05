<template>
  <header class="main-header">
    <img
      :src="
        isImageSet(conversationInfo.photo)
          ? serverUrl + '/images/' + conversationInfo.photo
          : getDefaultConversationImage(conversationInfo.category)
      "
      class="conversation-photo"
      alt="Foto conversazione"
    />
    <h1>{{ conversationInfo.name }}</h1>
    <div v-if="conversationInfo.category == 'group'" style="margin-left: auto">
      <input
        v-model="newMemberName"
        type="text"
        placeholder="Aggiungi membro..."
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
  <ul class="messages-container">
    <li
      v-for="message in messages"
      :key="message.id"
      :class="message.senderId == authId ? 'message-right' : 'message-left'"
      :style="
        message.senderId == authId
          ? 'align-self: flex-end; margin-right: 0.5rem;'
          : 'align-self: flex-start; margin-left: 0.5rem;'
      "
      @click="
        showMessageOptions =
          showMessageOptions == message.id ? null : message.id;
        showEmotes = false;
      "
    >
      <span v-if="message.senderId == authId" class="message-timestamp">
        {{ formatTimestamp(message.timestamp) }}
      </span>
      <span class="message-container"
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
      <label for="imageInput"
        ><img
          src="/icons/image.svg"
          class="send-photo-button"
          alt="Invia Immagine"
      /></label>
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
.main-header {
  display: flex;
  height: 6rem;
  align-items: center;
  border-bottom: solid black;
  background-color: rgb(0, 150, 250);
}
.conversation-photo {
  width: 80px;
  height: 80px;
  margin-left: 1rem;
  border: 0.2rem solid black;
  border-radius: 50%;
  object-fit: cover;
}
.add-to-group-input {
  height: 1.5rem;
  border: 0.15rem solid black;
  outline: none;
  font-size: 0.8rem;
  color: black;
}
.add-to-group-button {
  height: 2rem;
  border: 0.2rem solid black;
  background-color: gainsboro;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
}
.add-to-group-button:hover {
  background-color: grey;
}
.leave-group-button {
  height: 2rem;
  margin-left: 1rem;
  border: solid black;
  background-color: red;
  color: black;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
}
.leave-group-button:hover {
  background-color: darkred;
}
.messages-container {
  display: flex;
  flex-direction: column-reverse;
  overflow-y: auto;
  height: calc(100% - 9.5rem);
  margin: 0;
  padding: 0;
  list-style-type: none;
  background-color: gainsboro;
}
.message-right,
.message-left {
  position: relative;
  max-width: 60%;
  margin: 1rem 0rem;
  word-break: break-word;
  cursor: pointer;
}
.message-timestamp {
  margin-right: 0.5rem;
  font-size: 0.8rem;
}
.sender-name {
  margin-left: 0.5rem;
  white-space: nowrap;
  font-weight: bold;
  font-size: 0.95em;
}
.message-container {
  position: relative;
  padding: 0.5rem 0.7rem 0.5rem 0.7rem;
  background-color: rgb(0, 150, 250);
}
.message-container:hover {
  background-color: rgb(0, 100, 200);
}
.options {
  position: absolute;
  bottom: 100%;
  z-index: 1;
  white-space: nowrap;
}
.options-button {
  margin-right: 0.5rem;
  border: solid;
  border-color: black;
  border-radius: 0;
  font-size: 1rem;
  cursor: pointer;
}
.options-button-delete {
  margin-right: 0.5rem;
  border: solid;
  border-color: darkred;
  border-radius: 0;
  font-size: 1rem;
  color: darkred;
  cursor: pointer;
}
.emoji {
  position: absolute;
  top: 0.9rem;
  font-size: 1rem;
}
.forwarded {
  position: absolute;
  bottom: 100%;
  right: 0;
  padding: 0.1rem 0.4rem;
  background-color: rgb(0, 150, 250);
  font-size: 0.9rem;
  border-radius: 25%;
  font-weight: lighter;
}
.send-footer {
  position: relative;
  padding: 0.2rem 0.2rem 0 0.2rem;
  background-color: gainsboro;
}
.send-bar-wrapper {
  position: relative;
}
.send-bar {
  width: 100%;
  height: 3rem;
  box-sizing: border-box;
  border: 0.1rem solid black;
  outline: none;
  background-color: whitesmoke;
  color: black;
  font-size: 1.2rem;
}
.send-bar::placeholder {
  font-size: 1.2rem;
  color: black;
  font-style: italic;
}
.send-photo-button {
  position: absolute;
  right: 1rem;
  top: 0.5rem;
  z-index: 1;
  width: 1.6rem;
  height: 2rem;
  font-size: 1rem;
  cursor: pointer;
}
.check-icon {
  width: 0.9rem;
  height: 0.9rem;
  stroke: black;
  stroke-width: 2;
  fill: none;
  vertical-align: sub;
}
.settings-icon {
  width: 2rem;
  height: 2rem;
  margin-bottom: 3rem;
  stroke: black;
  stroke-width: 2;
  cursor: pointer;
}
</style>
