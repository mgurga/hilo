export { matchers } from './client-matchers.js';

export const components = [
	() => import("../runtime/components/layout.svelte"),
	() => import("../runtime/components/error.svelte"),
	() => import("../../src/routes/index.svelte"),
	() => import("../../src/routes/mygames.svelte"),
	() => import("../../src/routes/editor/[id].svelte"),
	() => import("../../src/routes/login.svelte")
];

export const dictionary = {
	"": [[0, 2], [1]],
	"mygames": [[0, 3], [1]],
	"editor/[id]": [[0, 4], [1]],
	"login": [[0, 5], [1]]
};