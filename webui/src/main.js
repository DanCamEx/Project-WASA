import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import axios from "./services/axios.js";
import LineIcons from "@lineiconshq/vue-lineicons";

const app = createApp(App);
app.config.globalProperties.$axios = axios;
app.use(router);
app.use(LineIcons);
app.mount("#app");
