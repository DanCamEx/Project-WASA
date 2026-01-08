import { createRouter, createWebHashHistory } from "vue-router";
import { userAuthId } from "../services/auth.js";
import LoginView from "../views/LoginView.vue";
import HomeView from "../views/HomeView.vue";
import ConversationView from "../views/ConversationView.vue";
import ForwardView from "../views/ForwardView.vue";
import WelcomeView from "../views/WelcomeView.vue";
import UserSettingsView from "../views/UserSettingsView.vue";
import GroupCreationView from "../views/GroupCreationView.vue";
import GroupSettingsView from "../views/GroupSettingsView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: "/", redirect: "/login" },
		{ path: "/login", name: "login", components: { default: LoginView } },
		{ path: "/home", components: { left: HomeView, right: WelcomeView } },
		{
			path: "/settings",
			components: { left: HomeView, right: UserSettingsView },
		},
		{
			path: "/conversations/:conversationId",
			components: {
				left: HomeView,
				right: ConversationView,
			},
		},
		{
			path: "/conversations/:conversationId/messages/:messageId/forward",
			components: {
				left: ForwardView,
				right: ConversationView,
			},
		},
		{
			path: "/new-group",
			components: {
				left: HomeView,
				right: GroupCreationView,
			},
		},
		{
			path: "/conversations/:conversationId/settings",
			components: {
				left: HomeView,
				right: GroupSettingsView,
			},
		}
	],
});

router.beforeEach((to, from, next) => {
  if (to.path != "/login" && !userAuthId.value) {
    next("/login");
  } else {
    next();
  }
});

export default router;
