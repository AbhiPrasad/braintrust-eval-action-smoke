import { Eval } from "braintrust";

function exactMatch({ output, expected }) {
  return {
    name: "exact_match",
    score: output === expected ? 1 : 0,
  };
}

Eval("Smoke JS Eval Action", {
  data: () => [
    { input: "JS", expected: "Hello JS" },
    { input: "GitHub Actions", expected: "Hello GitHub Actions" },
  ],
  task: async input => `Hello ${input}`,
  scores: [exactMatch],
});
