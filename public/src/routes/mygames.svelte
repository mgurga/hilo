<script lang="ts">
    import { username, key, server_url } from '../stores.js';
    import Nav from "../components/Nav.svelte";
    import { onMount } from 'svelte';

    type Game = {
        creator: string;
        id: string;
        name: string;
        description: string;
        datecreated: number;
    }

    let usergames: Array<Game> = [];

    onMount(() => {
        if ($username == "" || $key == "") {
            alert("you are not signed in");
            window.location.href = "/login";
        }

        fetch(`${$server_url}/api/usergames/${$username}`, {"method": "GET"})
        .then((response) => response.json())
        .then((data) => {
            usergames = data;
        })
    })

    function createGame() {
        fetch(`${$server_url}/api/creategame`, {"method": "GET", "headers" : { "authkey": $key }})
        .then((response) => response.json())
        .then((data) => {
            usergames = [...usergames, data]
        })
    }

    function deleteGame(id: string, index: number) {
        usergames = [...usergames.slice(0, index), ...usergames.slice(index + 1)];
        fetch(`${$server_url}/api/deletegame`, {"method": "GET", "headers" : { "authkey": $key, "id": id }})
    }
</script>

<Nav />
<main>
    <div id="menu">
        <div id="gamesheader">
            <h2 class="header">My Games</h2>
            <div style="float: right;">
                <div on:click={createGame} id="creategame">+</div>
            </div>
        </div>
        <div id="usergames">
            {#each usergames as g, i}
            <div class="gamecard">
                <div style="width: fit-content; display: inline-block;">
                    <h3 class="gamecardinfo">{g.name}</h3>
                    <br>
                    <br>
                    <p class="gamecardinfo">{g.description}</p>
                </div>
                <div style="float: right;">
                    <a style="text-decoration: none;" href="/editor/{g.id}">
                        <div class="divbutton" style="background-color: greenyellow;">Edit</div>
                    </a>
                    <a style="text-decoration: none;" href="/mygames" on:click={() => {deleteGame(g.id, i)}}>
                        <div class="divbutton" style="background-color: red;">Delete</div>
                    </a>
                </div>
            </div>
            {/each}
        </div>
    </div>
</main>

<style>
    .gamecardinfo {
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        display: inline;
        margin: 0;
        width: fit-content;
    }

    #usergames {
        width: 100%;
    }

    .gamecard {
        display: inline-block;
        width: 97%;
        padding: 10px;
        border: 1px solid black;
        border-top: 0;
    }

    .divbutton {
        padding-left: 10px;
        padding-right: 10px;
        margin: 5px;
        cursor: pointer;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        color: black;
    }

    #gamesheader {
        border-bottom: 2px solid black;
        padding: 0 10px 20px 10px;
    }

    .header {
        margin: 0;
        width: fit-content;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        display: inline;
        font-size: 30px;
    }

    #creategame {
        background-color: green;
        font-size: 30px;
        font-weight: bolder;
        padding-left: 10px;
        padding-right: 10px;
        margin: 5px;
        cursor: pointer;
    }

    main {
        display: flex;
        justify-content: center;
        /* align-items: center; */
        height: 95%; /* minus 5% to account for Nav bar */
    }

    #menu {
        margin-top: 20px;
        padding: 10px;
        width: 70%;
        /* border: 1px solid grey; */
        /* height: 50%; */
    }
</style>