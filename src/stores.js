import { writable } from 'svelte/store';
import { browser } from "$app/env";

let isProduction = import.meta.env.MODE === 'production';
console.log(import.meta.env.MODE);
export const server_url = isProduction ? writable("https://hiloserver.pythonanywhere.com") : writable("http://localhost:5000");

export let username = writable("");
export let key = writable("");

if (browser) {
    username.set(JSON.parse(localStorage.getItem("user")) || "")
    key.set(JSON.parse(localStorage.getItem("key")) || "")
    username.subscribe((value) => localStorage.setItem("user", JSON.stringify(value)))
    key.subscribe((value) => localStorage.setItem("key", JSON.stringify(value)))
}