# braintrust-eval-action-smoke

Disposable consumer repo for testing an unreleased `braintrustdata/eval-action` branch/SHA before cutting a release.

It runs three smoke jobs:

- `js-npm`: real Braintrust JS eval through `braintrust eval --jsonl`
- `python-pip`: real Braintrust Python eval through `braintrust eval --jsonl`
- `go-jsonl`: synthetic Go `go run .` program that emits a Braintrust-style experiment summary JSON line, which tests the action's Go runtime/parser/comment path

## Create and push the smoke repo

```bash
cd ~/workspace/braintrust-eval-action-smoke
git init
git add .
git commit -m "Add eval-action smoke tests"

gh repo create braintrust-eval-action-smoke --private --source . --remote origin --push
```

## Configure secrets/variables

Required secret:

```bash
gh secret set BRAINTRUST_API_KEY --body "$BRAINTRUST_API_KEY"
```

For `push` and `pull_request` runs, configure the eval-action branch/SHA under test as repo variables:

```bash
gh variable set EVAL_ACTION_REPOSITORY --body braintrustdata/eval-action
gh variable set EVAL_ACTION_REF --body <branch-or-sha-to-test>
```

For `workflow_dispatch` runs, you can pass these as inputs instead.

## Run manually

```bash
gh workflow run eval-action-smoke.yml \
  -f action_repository=braintrustdata/eval-action \
  -f action_ref=<branch-or-sha-to-test>

gh run watch
```

## Verify PR comments

The action only creates/updates PR comments when GitHub can associate the run with a PR. The easiest check is to open a PR from a branch in this same smoke repo after setting `EVAL_ACTION_REF`.

```bash
git checkout -b smoke/test-eval-action
# make a tiny change, e.g. edit this README
git commit -am "Trigger smoke workflow"
git push -u origin smoke/test-eval-action
gh pr create --fill
```

You should see one Braintrust eval report comment per `step_key`:

- `smoke-js-npm`
- `smoke-python-pip`
- `smoke-go-jsonl`

## Notes

- The eval-action branch/SHA you test must have `eval/dist/index.js` built and committed.
- PR secrets are available for PRs from branches in the same repo. Avoid testing secret-backed evals from forks unless you intentionally handle that case.
- If you also want package-manager coverage, add duplicate jobs for Node `pnpm` and Python `uv`.
