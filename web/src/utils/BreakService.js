import axios from 'axios';

const SERVER_URL = process.env.API_URL

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
    return axios.post(SERVER_URL, {'username': username, 'password': password})
                .then(response => response.data)
}

export function refreshWithToken(token) {
    return axios.post(SERVER_URL, {'token': token})
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
