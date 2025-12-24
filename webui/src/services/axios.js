import axios from "axios";

const instance = axios.create({
	baseURL: "http://localhost:3000",
	timeout: 1000 * 5,
});

instance.interceptors.request.use((config) => {
	if (!config.url.endsWith("/login")) {
		const token = localStorage.getItem("authId");
		if (token) {
			config.headers["Authorization"] = "Bearer " + token;
		}
	}
	return config;
});

export default instance;
