import axios from 'axios';

const SERVER_URL = {
    '_': 'https://func.skmobi.com/function/break',
    'dev': 'https://func.skmobi.com/function/break-dev',
    'local': 'http://localhost:9999',
}

var days = [];

var token = null;

export function setToken(newToken) {
    token = newToken
    if (typeof window !== "undefined") {
        if (newToken == null) {
            localStorage.removeItem("token")
        }
        else {
            localStorage.setItem("token", newToken)
        }
    }
}

export function getToken() {
	if (typeof window !== "undefined") {
		token = localStorage.getItem("token")
	}
	return token
}
export function isLoggedIn() {
    return getToken()
}

export function refreshWithLogin(username, password) {
    return axios.post(getServerEndpoint(), {'username': username, 'password': password})
                .then(response => response.data)
}

export function refreshWithToken(token) {
    return axios.post(getServerEndpoint(), {'token': token})
                .then(response => response.data)
}

export function updateData(data) {
    data.History = data.History.reverse()
    if (typeof window !== "undefined") {
    	localStorage.setItem("data", JSON.stringify(data))
    }
}

export function getData() {
    var data = null
	if (typeof window !== "undefined") {
        data = JSON.parse(localStorage.getItem("data"))
    }
    return data || {'History': []}
}

function getServerEndpoint() {
    var c = '_'
    if (typeof window !== "undefined") {
        c = localStorage.getItem("apiServer")
	}
    return SERVER_URL[c] || SERVER_URL['_']
}

export function setServerEndpoint(server) {
    if (typeof window !== "undefined") {
    	localStorage.setItem("apiServer", server)
    }
}
