<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();

    export let target: number = 0;
    export let num: number = 0;
    export let inc: boolean = true;

    function increment() {
        if (inc) {
            const animationDuration = 2000;
            const frameDuration = 1000 / 60;
            const totalFrames = Math.round( animationDuration / frameDuration );

            let frame = 0;
            const easeOutQuad = t => t * ( 2 - t );
            const counter = setInterval(() => {
                frame++;
                const progress = easeOutQuad( frame / totalFrames );
                const currentCount = Math.round( target * progress );

                num = currentCount;

                // If we’ve reached our last frame, stop the animation
                if ( frame === totalFrames ) {
                    clearInterval( counter );
                    dispatch("finished")
                }
            }, frameDuration);
        }
    }

    setTimeout(function() {
        increment()
    }, 100);
</script>

<h1>{num}</h1>