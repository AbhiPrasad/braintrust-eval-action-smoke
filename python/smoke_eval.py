from braintrust import Eval
from autoevals import Levenshtein


Eval(
    "Smoke Python Eval Action",
    data=lambda: [
        {"input": "Python", "expected": "Hello Python"},
        {"input": "GitHub Actions", "expected": "Hello GitHub Actions"},
        # Intentional mismatch to exercise eval-action regression reporting.
        {"input": "Braintrust", "expected": "Goodbye Braintrust"},
    ],
    task=lambda input: f"Hello {input}",
    scores=[Levenshtein],
)
