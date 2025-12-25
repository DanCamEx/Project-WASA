<template>
  <h1>Entra in WASAText!</h1>
  <input
    type="text"
    v-model="testo"
    placeholder="Inserisci il tuo nome utente..."
  />
  <button @click="doLogin">Accedi</button>
</template>

<script>
import axios from "../services/axios.js";

export default {
  data() {
    return {
      testo: "",
    };
  },
  methods: {
    async doLogin() {
      try {
        const response = await axios.post("/session", {
          name: this.testo,
        });
        const id = response.data.id;
        localStorage.setItem("authId", id);
        alert("Login effettuato!");
        this.$router.push("/home");
      } catch (error) {
        alert("Errore di login");
      }
    },
  },
};
</script>
