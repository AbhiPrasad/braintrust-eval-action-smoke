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
    // Intentional mismatch to exercise eval-action regression reporting.
    { input: "Braintrust", expected: "Goodbye Braintrust" },
  ],
  task: async input => `Hello ${input}`,
  scores: [exactMatch],
});
