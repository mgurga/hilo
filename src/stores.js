import { writable } from 'svelte/store';
import { browser } from "$app/environment";

let isProduction = import.meta.env.MODE === 'production';
export const server_url = writable(import.meta.env.VITE_SERVER_URL);

export let username = writable("");
export let key = writable("");

if (browser) {
    username.set(JSON.parse(localStorage.getItem("user")) || "")
    key.set(JSON.parse(localStorage.getItem("key")) || "")
    username.subscribe((value) => localStorage.setItem("user", JSON.stringify(value)))
    key.subscribe((value) => localStorage.setItem("key", JSON.stringify(value)))
}