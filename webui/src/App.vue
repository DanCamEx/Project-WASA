<template>
  <div class="main-container">
    <div class="left-panel">
      <h1 class="title">WASAText</h1>
      <RouterView name="left" />
    </div>
    <div class="right-panel">
      <RouterView name="right" />
    </div>
  </div>
</template>

<script setup>
import { RouterLink, RouterView, useRouter } from "vue-router";
import { ref } from "vue";

const isLoggedIn = ref(localStorage.getItem("authId"));
const router = useRouter();
</script>

<style>
body {
  margin: 0;
  padding: 0;
  height: 100%;
}
</style>

<style scoped>
.main-container {
  display: flex;
  height: 100vh;
}
.left-panel {
  flex: 1;
  border-right: 0.2rem solid black;
  overflow: hidden;
}
.right-panel {
  flex: 2;
}
.title {
  text-align: center;
  margin-right: 3rem;
  font-size: 1.6rem;
  font-family: sans-serif;
}
</style>
