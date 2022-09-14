<script lang="ts">
    import Dialog, { Title, Content, Actions } from "@smui/dialog";
    import { fade } from 'svelte/transition';
    import Button, { Label } from '@smui/button';
    import Fab, { Icon } from '@smui/fab';
    import { onMount } from 'svelte';
    import { server_url } from '../../../stores.js';
    import CountUp from './CountUp.svelte';
    import { page } from '$app/stores';
    import { base } from "$app/paths";

    export let gameid = $page.params.slug;

    type Game = {
        creator: string;
        id: string;
        name: string;
        description: string;
        datecreated: number;
    }

    type GameNode = {
        name: string;
        amount: number;
        id: string;
        parent: string;
    }

    let gameinfo: Game = null;
    let gamenodes: Array<GameNode> = [];
    let nodelist: Array<GameNode>;
    let curnode: number = 0;
    let buttonsvisible = true;
    let showamounts = false;
    let outcome = 0; // 0 = nothing, 1 = win, 2 = lose, 3 = restart
    let userchoice = "";
    let score = 0;
    let progress = 0;
    let colors = ["darkblue", "rgb(0, 0, 90)"];
    let mandatorydialog = false;
    let infodialog = false;

    onMount(() => {
        fetch(`${$server_url}/api/game/${gameid}/info`, {"method": "GET"})
        .then((response) => response.json())
        .then((data) => {
            if (data.error == "invalid game") {
                alert("invalid game");
                window.location.pathname = base;
            }
            gameinfo = data;
        })
        fetch(`${$server_url}/api/game/${gameid}/nodes`, {"method": "GET"})
        .then((response) => response.json())
        .then((data) => {
            gamenodes = data;

            if (gamenodes.length >= 3) {
                nodelist = [gamenodes[(Math.random() * gamenodes.length) | 0]];
                let newitem = gamenodes[(Math.random() * gamenodes.length) | 0];
                for (var i = 0; i < 1000; i++) {
                    while (newitem.id == nodelist[nodelist.length - 1].id) {
                        newitem = gamenodes[(Math.random() * gamenodes.length) | 0];
                    }
                    nodelist = [...nodelist, newitem]
                }

                curnode = 3;
                infodialog = true;

                // console.log(nodelist.slice(0, 10));
            } else {
                mandatorydialog = true;
                buttonsvisible = false;
            }
        })
    })

    function decisionpressed(dir: string) {
        buttonsvisible = false;
        userchoice = dir;
        setTimeout(function() {
            showamounts = true;
        }, 1000)
    }

    function finishedcount() {
        if (userchoice == "higher") {
            outcome = nodelist[curnode].amount >= nodelist[curnode - 1].amount ? 1 : 2;
        } else {
            outcome = nodelist[curnode].amount <= nodelist[curnode - 1].amount ? 1 : 2;
        }

        if (outcome == 1) {
            score++;
        }
        
        setTimeout(function() {
            // showamounts = false;

            function increment() {
                setTimeout(function() {
                    if (outcome == 1 || outcome == 0) {
                        outcome = 0;
                        if (progress != 100) {
                            progress = progress + 1;
                            increment();
                        } else {
                            quickswitch();
                        }
                    } else {
                        outcome = 3;
                    }
                }, 10);
            }

            increment();
        }, 3000);
        // console.log(`${leftitem.amount} vs ${rightitem.amount} done counting, you ${outcome ? "WIN" : "LOSE"}`);
    }

    function quickswitch() {
        curnode++;
        colors = [colors[1], colors[0]];
        progress = 0;
        showamounts = false;
        buttonsvisible = true;
    }

    function restart() {
        if (outcome == 3) {
            nodelist = [gamenodes[(Math.random() * gamenodes.length) | 0]];
            let newitem = gamenodes[(Math.random() * gamenodes.length) | 0];
            for (var i = 0; i < 1000; i++) {
                while (newitem.id == nodelist[nodelist.length - 1].id) {
                    newitem = gamenodes[(Math.random() * gamenodes.length) | 0];
                }
                nodelist = [...nodelist, newitem];
            }

            score = 0;
            curnode = 3;
            buttonsvisible = true;
            showamounts = false;
            outcome = 0;
        }
    }
</script>
<Dialog
    bind:open={mandatorydialog}
    scrimClickAction=""
    escapeKeyAction="">
    <Title id="mandatory-title">Error</Title>
    <Content id="mandatory-content">The game does not have enough items to play.</Content>
</Dialog>
<Dialog
    bind:open={infodialog}>
    <Title>{gameinfo == null ? "" : gameinfo.name} by {gameinfo == null ? "" : gameinfo.creator}</Title>
    <Content>{gameinfo == null ? "" : gameinfo.description}</Content>
    <Actions>
        <Button on:click={() => (infodialog = false)}><Label>Play</Label></Button>
    </Actions>
</Dialog>

<title>HiLo</title>
<main style="background-color: {colors[0]}">
    <div id="left" style="background-color: {colors[0]}">
        <h1>{nodelist == null ? "" : nodelist[curnode - 1].name}</h1>
        {#if nodelist != null}
        <div class="amounts" transition:fade>
            <CountUp bind:num={nodelist[curnode - 1].amount} inc={false}/>
        </div>
        {/if}
    </div>
    <div id="right" style="transform: translateX(-{progress}%); background-color: {colors[1]};">
        <h1>{nodelist == null ? "" : nodelist[curnode].name}</h1>
        {#if buttonsvisible}
        <div id="buttons" style="display: grid;" transition:fade>
            <Button variant="raised" style="background-color: green; margin-bottom: 10px" on:click={() => {decisionpressed("higher")}}>
                <Label>Higher</Label>
                <Icon class="material-icons">arrow_upward</Icon>
            </Button>
            <Button variant="raised" style="background-color: red;" on:click={() => {decisionpressed("lower")}}>
                <Label>Lower</Label>
                <Icon class="material-icons">arrow_downward</Icon>
            </Button>
        </div>
        {/if}
        {#if showamounts}
        <div class="amounts" in:fade>
            <CountUp on:finished={finishedcount} target={nodelist[curnode].amount} />
        </div>
        {/if}
    </div>
    <div id="outcome">
        {#if outcome != 0}
        <div transition:fade>
            <Fab style="background-color: {outcome == 1 ? "green" : outcome == 2 ? "red" : "orange"}" on:click={restart}>
                <Icon class="material-icons">{outcome == 1 ? "check" : outcome == 2 ? "close" : "loop"}</Icon>
            </Fab>
        </div>
        {/if}
    </div>
    <h1 style="z-index: 6; position: absolute; color: white; margin-top: 5px">{score}</h1>
</main>

<style>
    h1 {
        font-size: 50px;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        justify-content: center;
    }

    #outcome {
        position: absolute;
        z-index: 5;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
        top: calc((100% / 2) - 5%);
    }

    #left, #right {
        width: 50%;
        margin: 0;
        height: 100%;
        color: white;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
    }

    #left > * {
        display: flex;
    }

    main {
        display: flex;
        justify-content: center;
        height: 100%; /* minus 5% to account for Nav bar */
    }
</style>