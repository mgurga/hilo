<script>
    import { onMount } from 'svelte';
    onMount(() => {
        window.location.href = import.meta.env.VITE_WEBSITE_BASE_URL + "/";
    })
</script>