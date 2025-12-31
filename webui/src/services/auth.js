import { ref } from "vue";

export const userAuthId = ref(sessionStorage.getItem("userAuthId"));
export const userName = ref(sessionStorage.getItem("userName"));
export const userPhoto = ref(sessionStorage.getItem("userPhoto"));
