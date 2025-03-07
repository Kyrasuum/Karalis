import axios from 'axios';
import { parallel } from 'src/types/net';
import store from 'src/store';

let base_url = '';
if(import.meta.env.PROD) {
	base_url = import.meta.env.VITE_ASSET_URL;
} else {
	base_url = 'http://localhost:3000' + import.meta.env.VITE_ASSET_URL;
}
const url = base_url
const api = axios.create({ baseURL: url });

async function ApiQuery(endpoint: string, req: any, limit = 0): Promise<string> {
	let data = null;

	let attempt = 1;
	while (data == null) {
		data = await api.post<string>(endpoint, req)
		.then((response) => {
			return response.data;
		}).catch(() => {
			return "{\"payload\":\"error";
		});

		if (limit > 0 && attempt >= limit) {
		    break;
		}
		if (IsNetError(data)) {
			data = null;
			await new Promise((f) => setTimeout(f, 100 * attempt));
			if (attempt < 50) {
				attempt*=2;
			} else {
				attempt++;
			}
		}
	}

	return data;
};
const api_raw = {
	//raw functions for sending and recieving api requests
};

export { api_raw, url };
