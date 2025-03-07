import router from 'src/router';
import {setCookie, getCookie, url_enc, kc_enc, sess_enc} from 'src/types/Cookie';
import {getEncoding, getEndpoint} from 'src/types/Path';
import {api} from 'src/api/api';
import store from 'src/store';

let config_default = {
    "url": "",
    "realm": import.meta.env.VITE_REALM,
    "clientId": import.meta.env.VITE_CLIENT,
    'ssl-required': 'external',
    'public-client': true
}

let init_default = {
    checkLoginIframe: false,
	enableLogging: "true",
	onLoad: "login-required",
	responseMode: "query",
	scope: import.meta.env.VITE_SCOPE,
	redirectUri: "", // override with public url
}

let keycloak_default = {
    idTokenParsed: {
        email: "user@gmail.com",
        family_name: "user",
        given_name: "user",
        preferred_username: "user@gmail.com",
        user_id: "test-test-test",
    }
};

export function Init() {
    let init = init_default;
    init.redirectUri = import.meta.env.VITE_REDIRECT_URL + import.meta.env.VITE_ASSET_URL + "/";
    let endpoint = getEndpoint();
    init.redirectUri += endpoint;

    const ret = init
    return ret;
}

export function Config() {
    let config = config_default;
    config["url"] = import.meta.env.VITE_AUTH_URL;

    let enc = getEncoding();
    if (enc.length > 0) {
        setCookie(url_enc, enc);
    }

    const ret = config
    return ret;
}

export function onReady(kc) {
    let endpoint = getEndpoint();
    let enc = getCookie(url_enc);
    if (enc.length > 0) {
	    router.push({ path: endpoint,  query: {'': enc} });
    }
}

export function onInitError(error) {
    console.log(error);
}
