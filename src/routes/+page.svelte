<script lang="ts">
    import Nav from "../components/Nav.svelte";
    import LayoutGrid, { Cell } from '@smui/layout-grid';
    import { fade } from 'svelte/transition';
    import { onMount } from 'svelte';
    import { server_url } from '../stores.js';

    type Game = {
        creator: string;
        id: string;
        name: string;
        description: string;
        datecreated: number;
    }

    let recents: Array<Game> = null;

    onMount(() => {
        fetch(`${$server_url}/api/recent/3`, {"method": "GET"})
        .then((response) => response.json())
        .then((data) => {
            recents = data;
        })
    })
</script>

<title>HiLo</title>
<Nav />
<main>
    <div id="menu">
        <div id="header">
            <h1>HiLo: Higher or Lower</h1>
            <p>Play and Create you own higher or lower games</p>
        </div>

        {#if recents !== null && recents.length > 0}
        <div id="recentgames" transition:fade>
            <h1 style="margin-top: 20%; margin-bottom: 0;">Recent Games:</h1>
            <LayoutGrid>
            {#each recents as game}
                <Cell>
                    <div class="gamecell" on:click={() => {window.location.href = import.meta.env.VITE_WEBSITE_BASE_URL + `/play/${game.id}`}}>
                        <h1 style="margin-bottom: 0; margin-top: 0;">{game.name}</h1>
                        <h3 style="color: grey; margin-bottom: 0;">{game.description}</h3>
                        <p style="color: grey; margin-bottom: 0;">created by: {game.creator}</p>
                    </div>
                </Cell>
            {/each}
            </LayoutGrid>
        </div>
        {/if}
    </div>
</main>

<style>
    main {
        display: flex;
        justify-content: center;
        align-items: center;
        height: calc(100% - 64px); /* minus 5% to account for Nav bar */
    }

    #menu {
        text-align: center;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    }

    .gamecell {
        width: 300px;
        height: 200px;
        display: flex;
        justify-content: center;
        align-items: center;
        background-color: var(--mdc-theme-secondary, #333);
        color: var(--mdc-theme-on-secondary, #fff);
        flex-direction: column;
        cursor: pointer;
    }
</style>