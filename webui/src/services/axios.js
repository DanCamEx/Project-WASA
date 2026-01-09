import axios from "axios";

const serverURL = __API_URL__;
const instance = axios.create({
	baseURL: serverURL,
	timeout: 1000 * 5,
});

instance.interceptors.request.use((config) => {
	if (!config.url.endsWith("/login")) {
		const token = sessionStorage.getItem("userAuthId");
		if (token) {
			config.headers["Authorization"] = "Bearer " + token;
		}
	}
	return config;
});

export default instance;
