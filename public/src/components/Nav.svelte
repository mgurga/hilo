<script lang="ts">
    import { username, key } from '../stores.js';
    import Logout from 'svelte-material-icons/Logout.svelte';

    function logout() {
        $username = "";
        $key = "";
    }
</script>

<nav>
    <a href="/"><span>HiLo</span></a>
    <div id="floatright">
        {#if $key == ""}
        <a href="/login"><span>Login</span></a>
        {:else}
        <a href="/mygames"><span>My Games</span></a>
        <span style="padding-right: 0px">Welcome, {$username}</span>
        <a on:click={logout} href="/login">
            <span style="padding: 2px; height: 100%;">
                <Logout color="white" height="90%"/>
            </span>
        </a>
        {/if}
    </div>
</nav>

<style>
    nav {
        height: 5%;
        background-color: rgb(65, 179, 218);
    }

    span {
        padding: 10px;
        display: inline-block;
        color: white;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    }

    a:hover > * {
        color: black;
        cursor: pointer;
        background-color: lightblue;
    }

    #floatright {
        float: right;
    }
</style>