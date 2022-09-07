<script lang="ts">
    import { username, key } from '../stores.js';
    import IconButton from '@smui/icon-button';
    import Button from '@smui/button';
    import TopAppBar, { Row, Section, Title } from '@smui/top-app-bar';
    import { base } from '$app/paths';

    function logout() {
        $username = "";
        $key = "";
        window.location.pathname = `${base}/login`
    }
</script>

<nav>
    <TopAppBar
      variant="static"
      color='primary'>
      <Row>
        <Section>
          <Title on:click={() => {window.location.pathname = base}} style="cursor: pointer; user-select: none;">HiLo</Title>
        </Section>
        <Section align="end" toolbar>
            {#if $key != ""}
            <Button on:click={() => {window.location.pathname = `${base}/mygames`}}>My Games</Button>
            <IconButton on:click={logout} class="material-icons" aria-label="Logout">logout</IconButton>
            {:else}
            <Button on:click={() => {window.location.pathname = `${base}/login`}}>Login</Button>
            {/if}
        </Section>
      </Row>
    </TopAppBar>
</nav>