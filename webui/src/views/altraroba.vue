<ul class="list-group mx-3">
          <li
            v-for="user in filteredUsers"
            :key="user.id"
            class="list-group-item"
          >
            <svg
              class="feather mb-1 rounded-circle border"
              style="width: 40px; height: 40px"
            >
              <use href="/feather-sprite-v4.29.0.svg#user" />
            </svg>
            {{ user.name }}
            <button
              class="btn btn-sm btn-dark float-end mt-1"
              @click="createConversation(user.id, 'direct')"
            >
              Crea conversazione
            </button>
          </li>
        </ul>
        <h5 class="fw-bold ms-4 mb-2 mt-3">Le tue conversazioni</h5>
        <ul class="list-group flex-grow-1">
          <li
            v-for="conversation in conversations"
            :key="conversation.id"
            class="list-group-item"
            @click="getConversation(conversation.id)"
            style="cursor: pointer"
          >
            <svg
              class="feather mb-1 rounded-circle border"
              style="width: 40px; height: 40px"
            >
              <use href="/feather-sprite-v4.29.0.svg#user" />
            </svg>
            <h8 class="fw-bold ms-2 me-3">{{ conversation.name }}</h8>
            ultimo messaggio:
            {{ latestMessageContents[conversation.id] }}
          </li>
        </ul>
  
      <div class="col-8 d-flex flex-column vh-100">
        <div
          class="bg-light border border-dark py-3 px-3 mt-2"
          v-if="selectedConversation"
        >
          {{ selectedConversation.name }}
        </div>
	  </div>
        <div class="flex-grow-1 d-flex flex-column-reverse overflow-auto">
          <div
            v-for="message in messages.slice().reverse()"
            :key="message.id"
            class="mb-2 fs-6"
            :class="
              Number(message.senderId) == Number(authId)
                ? 'text-end'
                : 'text-start'
            "
          >
		  </div>
            <div
              class="d-flex me-2"
              v-if="subMenuMessageId === message.id"
              :class="
                Number(message.senderId) == Number(authId)
                  ? 'justify-content-end'
                  : 'justify-content-start'
              "
            >
              <h5
                v-if="commentEmojiMenuVisible"
                style="cursor: pointer"
                @click="
                  commentMessage('👍');
                  commentEmojiMenuVisible = false;
                  subMenuMessageId = null;
                "
              >
                👍
              </h5>
              <h5
                v-if="commentEmojiMenuVisible"
                style="cursor: pointer"
                @click="
                  commentMessage('😂');
                  commentEmojiMenuVisible = false;
                  subMenuMessageId = null;
                "
              >
                😂
              </h5>
              <h5
                v-if="commentEmojiMenuVisible"
                style="cursor: pointer"
                @click="
                  commentMessage('💀');
                  commentEmojiMenuVisible = false;
                  subMenuMessageId = null;
                "
              >
                💀
              </h5>
              <h5
                v-if="commentEmojiMenuVisible"
                style="cursor: pointer"
                @click="
                  commentMessage('😭');
                  commentEmojiMenuVisible = false;
                  subMenuMessageId = null;
                "
              >
                😭
              </h5>
              <h5
                v-if="commentEmojiMenuVisible"
                style="cursor: pointer"
                @click="
                  commentMessage('🔥');
                  commentEmojiMenuVisible = false;
                  subMenuMessageId = null;
                "
              >
                🔥
              </h5>
              <h4
                class="btn btn-sm btn-dark me-2"
                @click="commentEmojiMenuVisible = !commentEmojiMenuVisible"
              >
                Commenta
              </h4>
              <h4 class="btn btn-sm btn-dark" @click="">Inoltra</h4>
            </div>
            <svg
              class="feather mt-1 ms-1"
              style="cursor: pointer"
              @click="
                subMenuMessageId =
                  message.id === subMenuMessageId ? null : message.id;
                commentEmojiMenuVisible = false;
              "
            >
              <use href="/feather-sprite-v4.29.0.svg#chevron-down" />
            </svg>
            <span
              class="d-inline-block p-2 border border-dark rounded"
              style="
                max-width: 70%;
                word-break: break-word;
                background: #f8f9fa;
              "
            >
              {{ message.content }}
              <svg
                v-if="message.senderId == authId && message.status == 'read'"
                class="feather ms-1"
              >
                <use href="/feather-sprite-v4.29.0.svg#check-circle" />
              </svg>
              <svg v-else-if="message.senderId == authId" class="feather ms-1">
                <use href="/feather-sprite-v4.29.0.svg#check" />
              </svg>
            </span>
            <div
              v-for="comment in messageComments[message.id] || []"
              :key="comment.reacterId + '-' + comment.messageId"
            >
              <div
                class="d-inline-block p-1 mb-1 border border-dark rounded"
                style="
                  background: #e9ecef;
                  max-width: 20%;
                  word-break: break-word;
                "
              >
                {{ comment.emoji }}
              </div>
            </div>
          </div>
    
        <div class="input-group mt-auto mb-3">
          <input
            type="text"
            class="form-control mt-auto mb-3"
            placeholder="Scrivi un messaggio..."
            :hidden="!selectedConversation"
            v-model="messageContent"
            @keyup.enter="sendMessage()"
          />
          <svg
            class="feather mt-1 ms-1"
            style="width: 30px; height: 30px; cursor: pointer"
            :hidden="!selectedConversation"
          >
            <use href="/feather-sprite-v4.29.0.svg#paperclip" />
          </svg>
		</div>

   
   
</template>

<script>
import axios from "../services/axios.js";
import UserSearchBar from "../components/UserSearchBar.vue";

export default {
  components: { UserSearchBar },
  data() {
    return {
      conversations: [],
      messages: [],
      filteredUsers: [],
      allUsers: [],
      messageComments: {},
      latestMessageContents: {},
      selectedConversation: null,
      messageContent: "",
      userSearchText: "",
      authId: localStorage.getItem("authId"),
      subMenuMessageId: null,
      commentEmojiMenuVisible: false,
      forwardingMessageId: null,
    };
  },
  methods: {
    async getMyConversations() {
      try {
        const response = await axios.get("/users/me/conversations");
        this.conversations = response.data;
        this.conversations.forEach((conversation) => {
          this.getLatestMessage(conversation.latestMessageId, conversation.id);
        });
      } catch (error) {
        alert("Errore nel caricamento delle conversazioni");
      }
    },
    async getConversation(id) {
      try {
        const response = await axios.get(`/conversations/${id}/messages`);
        this.messages = response.data;
        this.selectedConversation = this.conversations.find(
          (conversation) => Number(conversation.id) == Number(id)
        );
        for (const message of this.messages) {
          this.getComments(message.id);
        }
      } catch (error) {
        alert("Errore nel caricamento dei messaggi della conversazione");
      }
    },
    async sendMessage() {
      try {
        await axios.post(
          `/conversations/${this.selectedConversation.id}/messages`,
          {
            content: this.messageContent,
          }
        );
        this.messages.push({
          senderId: this.authId,
          content: this.messageContent,
          status: "sent",
        });
        this.latestMessageContents[this.selectedConversation.id] =
          this.messageContent;
        this.messageContent = "";
      } catch (error) {
        alert("Errore nell'invio del messaggio");
      }
    },
    async searchUsers() {
      try {
        const response = await axios.get("/users/search");
        this.allUsers = response.data;
      } catch (error) {
        alert("Errore nella ricerca degli utenti");
      }
    },
    filterUsers(query) {
      const filter = query.toLowerCase();
      this.filteredUsers = this.allUsers.filter((user) =>
        user.name.toLowerCase().includes(filter)
      );
      if (!filter) {
        this.filteredUsers = [];
      }
    },
    async createConversation(userId, category) {
      try {
        const response = await axios.post("/users/me/conversations", {
          participantIds: [userId],
          category: category,
        });
        this.conversations.push(response.data);
        alert("Conversazione creata con successo!");
      } catch (error) {
        alert("Errore nella creazione della conversazione");
      }
    },
    async getLatestMessage(id, conversationId) {
      if (!id) {
        this.latestMessageContents[conversationId] = "";
        return;
      }
      try {
        const response = await axios.get(`/messages/${id}`);
        this.latestMessageContents[conversationId] = response.data.content;
      } catch (error) {
        alert("Errore nel recupero del messaggio");
      }
    },
    async getComments(messageId) {
      try {
        const response = await axios.get(`/messages/${messageId}/comments`);
        this.messageComments[messageId] = response.data;
      } catch (error) {
        alert("Errore nel recupero dei commenti");
      }
    },
    async commentMessage(emoji) {
      try {
        await axios.post(`/messages/${this.subMenuMessageId}/comments`, {
          emoji: emoji,
        });
        alert("Commento aggiunto con successo!");
        if (!this.messageComments[this.subMenuMessageId]) {
          this.messageComments[this.subMenuMessageId] = [];
        }
        console.log(
          "Prima del push:",
          this.messageComments[this.subMenuMessageId]
        );
        this.messageComments[this.subMenuMessageId].push({
          reacterId: this.authId,
          messageId: this.subMenuMessageId,
          emoji: emoji,
        });
        console.log(
          "Dopo il push:",
          this.messageComments[this.subMenuMessageId]
        );
        console.log("messageComments globale:", this.messageComments);
      } catch (error) {
        alert("Errore nell'aggiunta del commento");
      }
    },
    async uncommentMessage(messageId) {
      try {
        await axios.delete(`/messages/${messageId}/comments`);
        alert("Commento rimosso con successo!");
        this.messageComments[messageId] = (
          this.messageComments[messageId] || []
        ).filter((comment) => comment.messageId !== messageId);
      } catch (error) {
        alert("Errore nella rimozione del commento");
      }
    },
    async forwardMessage(conversationId) {
      try {
        const message = this.messages.find(
          (m) => m.id === this.forwardingMessageId
        );
        if (!message) return;
        await axios.post(`/conversations/${conversationId}/messages`, {
          content: message.content,
        });
        alert("Messaggio inoltrato!");
        this.forwardingMessageId = null;
      } catch (error) {
        alert("Errore nell'inoltro del messaggio");
        this.forwardingMessageId = null;
      }
    },
    cancelForward() {
      this.forwardingMessageId = null;
    },
  },
  mounted() {
    this.getMyConversations();
    this.searchUsers();
  },
};
</script>
