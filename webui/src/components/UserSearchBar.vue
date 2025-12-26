<template>
  <input
    type="text"
    placeholder="Cerca un utente..."
    v-model="searchQuery"
    class="search-bar"
    @input="onSearch"
  />
</template>

<script setup>
import { ref } from "vue";
const searchQuery = ref("");
const emit = defineEmits(["search"]);
const onSearch = () => {
  emit("search", searchQuery.value);
};
</script>

<style scoped>
.search-bar {
  width: 95%;
  height: 2rem;
  margin-left: 0.8rem;
  margin-top: 1rem;
  border: 0.1rem solid #000000;
  font-family: sans-serif;
  font-size: 1rem;
  font-style: italic;
  background-color: #f1f1f1;
  color: #000000;
  outline: none;
}
.search-bar::placeholder {
  color: #000000;
}
</style>
