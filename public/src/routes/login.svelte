<script lang="ts">
    import Nav from "../components/Nav.svelte";
    import { username, key, server_url } from '../stores.js';
    import { onMount } from 'svelte';

    // console.log("server url: " + $server_url);
    enum Menu {
        Register,
        Login
    }
    
    let cur: Menu = Menu.Login;
    let error = "";
    let success = "";

    onMount(() => {
        const url = window.location.href
        if (url.includes("#")) {
            let hash = url.split("#")[1];
            let user = url.split("#")[2];
            if (hash == "incorrect") {
                error = "user does not exist";
            } else if(hash == "exists") {
                error = "username is already registered";
            } else if(hash == "registered") {
                success = "successfully registered, now login";
            } else {
                fetch(`${$server_url}/api/authkey/${hash}/${user}`, {"method": "GET"})
                .then((response) => response.text())
                .then((data) => {
                    if (data.length == 20) {
                        $key = data;
                        $username = user;
                        success = "successfully logged in";
                        // console.log(data);
                        window.history.pushState({}, document.title, "/" + "login");
                    }
                    if (data == "")
                        alert("invalid hash, try signing in again")
                })
            }
        }
    });
</script>

<title>Login</title>
<Nav />
<main>
    <div id="menu">
        <button on:click={() => cur = Menu.Login} style="width: 48%">Login</button>
        <button on:click={() => cur = Menu.Register} style="width: 48%">Register</button>
        <form action="{$server_url}/api/{cur == Menu.Login ? "signin" : "register"}">
            <h2 style="width: auto; text-align: center;">{cur == Menu.Login ? "Login" : "Register"}</h2>
            <p>Username</p>
            <input type="text" name="user">
            <br>
            <br>
            <p>Password</p>
            <input type="password" name="pass">
            <br>
            <br>
            <button type="submit">Submit</button>
        </form>
        <p style="color: red;">{error}</p>
        <p style="color: green;">{success}</p>
    </div>
</main>

<style>
    main {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 95%; /* minus 5% to account for Nav bar */
    }

    h2, p {
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    }

    p {
        margin: 0;
        color: rgb(43, 43, 43);
    }

    #menu {
        padding: 10px;
        border: 1px solid grey;
        width: 30%;
        /* height: 50%; */
    }
</style>