from braintrust import Eval
from autoevals import Levenshtein


Eval(
    "Smoke Python Eval Action",
    data=lambda: [
        {"input": "Python", "expected": "Hello Python"},
        {"input": "GitHub Actions", "expected": "Hello GitHub Actions"},
    ],
    task=lambda input: f"Hello {input}",
    scores=[Levenshtein],
)
