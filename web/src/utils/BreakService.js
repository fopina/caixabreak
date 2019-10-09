import axios from 'axios';

const SERVER_URL = process.env.PREACT_APP_API_URL || 'http://localhost:9999'

export function setToken(token) {
    if (typeof window !== "undefined") {
        if (token == null) {
            localStorage.removeItem("token")
        }
        else {
            localStorage.setItem("token", token)
        }
    }
}

export function getToken() {
	if (typeof window !== "undefined") {
		return localStorage.getItem("token")
	}
	return ""
}

export function saveLogin(username, password) {
    if (typeof window !== "undefined") {
        localStorage.setItem("u", username)
        localStorage.setItem("p", password)
    }
}

export function removeLogin() {
    if (typeof window !== "undefined") {
        localStorage.removeItem("u")
        localStorage.removeItem("p")
    }
}

export function isLoggedIn() {
    return getToken()
}

export function refreshWithLogin(username, password) {
    return axios.post(SERVER_URL, {'username': username, 'password': password})
                .then(response => response.data)
}

export function refreshWithSavedLogin() {
    return refreshWithLogin(localStorage.getItem("u"), localStorage.getItem("p"))
}

export function refreshWithToken() {
    return axios.post(SERVER_URL, { 'token': getToken() })
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
