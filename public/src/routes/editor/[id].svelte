<script lang="ts">
    import { username, key, server_url } from '../../stores.js';
    import Nav from "../../components/Nav.svelte";
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import Dialog, { Title, Actions } from "@smui/dialog";
    import Button, { Label } from '@smui/button';
    import Textfield from '@smui/textfield';
    import CircularProgress from '@smui/circular-progress';
    // console.log("editing: " + $page.params.id)

    type Game = {
        creator: string;
        id: string;
        name: string;
        description: string;
        datecreated: number;
    }

    type Node = {
        name: string;
        amount: number;
        id: string;
    }

    let gameinfo: Game = null;
    let gamenodes: Array<Node> = [];
    let open = false;
    let newName = "";
    let newDesc = "";

    onMount(() => {
        fetch(`${$server_url}/api/game/${$page.params.id}`, {"method": "GET"})
        .then((response) => response.json())
        .then((data) => {
            gameinfo = data;
        })
    })

    function changename() {
        newName = gameinfo.name;
        newDesc = gameinfo.description;
        open = true;
    }

    function savename() {
        open = false;
        verifylogin();
        fetch(`${$server_url}/api/setnamedesc`, 
            {"method": "POST", "headers": {"authkey": $key, "id": $page.params.id, "name": newName, "desc": newDesc}})
            .then(() => {
                gameinfo.name = newName;
                gameinfo.description = newDesc;
            })
    }

    function verifylogin() {
        if ($username == "" || $key == "") {
            alert("you are not signed in")
            window.location.href = "/login";
        }
    }
</script>

<Nav />
<main>
    <div id="menu">
        <div id="gamebanner">
            {#if gameinfo == null}
            <div style="display: flex; justify-content: center">
                <CircularProgress style="height: 32px; width: 32px;" indeterminate />
            </div>
            {:else}
            <h1 on:click={changename} class="gameinfotext">{gameinfo.name}</h1>
            <br>
            <br>
            <p on:click={changename} class="gameinfotext">{gameinfo.description}</p>
            <div style="float: right; display: inline;">
                <p>created by: {gameinfo.creator}</p>
            </div>
            {/if}
        </div>
    </div>
</main>
<Dialog
    bind:open
    aria-labelledby="simple-title"
    aria-describedby="simple-content">
    <Title id="simple-title">Rename</Title>
    <Textfield bind:value={newName} label="Name" />
    <br>
    <Textfield textarea bind:value={newDesc} label="Description" />
    <Actions>
        <Button on:click={savename}>
        <Label>Save</Label>
        </Button>
    </Actions>
</Dialog>

<style>
    .gameinfotext {
        margin: 0;
        cursor: pointer;
        user-select: none;
        width: fit-content;
        display: inline-block;
    }

    #gamebanner {
        background-color: white;
        color: black;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        padding: 10px 10px 15px 10px;
        border-radius: 5px;
        border: 1px solid black;
    }

    main {
        display: flex;
        justify-content: center;
        height: 95%; /* minus 5% to account for Nav bar */
    }

    #menu {
        margin-top: 0;
        padding: 10px;
        width: 70%;
    }
</style>