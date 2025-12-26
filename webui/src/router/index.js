import { createRouter, createWebHashHistory } from "vue-router";
import LoginPage from "../views/LoginPage.vue";
import HomePage from "../views/HomePage.vue";
import ConversationPage from "../views/ConversationPage.vue";
import ForwardPage from "../views/ForwardPage.vue";
import WelcomePage from "../views/WelcomePage.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: "/", redirect: "/login" },
		{ path: "/login", name: "login", components: { default: LoginPage } },
		{ path: "/home", components: { left: HomePage, right: WelcomePage } },
		{
			path: "/conversations/:id",
			components: {
				left: HomePage,
				right: ConversationPage,
			},
		},
		{
			path: "/messages/:id/forward",
			components: {
				left: ForwardPage,
				right: ConversationPage,
			},
		},
	],
});

export default router;
