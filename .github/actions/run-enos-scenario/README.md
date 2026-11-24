# run-enos-scenario

Reusable composite action for running an Enos scenario with standardized retry, debug artifact upload, and destroy cleanup behavior.

## Behavior

The action performs this lifecycle:

1. Optionally installs Terraform.
2. Optionally installs Enos.
3. Prepares a sanitized debug artifact name from the scenario filter.
4. Runs `enos scenario launch`.
5. Retries the launch command once if the first attempt fails.
6. Optionally prints `enos scenario exec --cmd show` and `plan` before retry and after a failed retry.
7. Uploads debug data on failure.
8. Always destroys the scenario and retries destroy once if needed.

## Inputs

| Name | Required | Default | Description |
| --- | --- | --- | --- |
| `scenario-filter` | yes | n/a | Scenario filter passed to Enos. |
| `working-directory` | no | `./enos` | Value passed to `--chdir`. |
| `timeout` | no | `45m0s` | Timeout for the first main command attempt. |
| `retry-timeout` | no | `30m0s` | Timeout for the retry attempt. |
| `destroy-timeout` | no | `10m0s` | Timeout for destroy. |
| `show-state-on-retry` | no | `false` | When `true`, prints scenario state and plan before retry and after failed retry. |
| `upload-debug-data` | no | `true` | When `true`, uploads debug data on failure. |
| `debug-data-retention-days` | no | `30` | Artifact retention days for debug data. |
| `debug-data-path` | no | empty | Explicit debug data directory. If empty, uses `ENOS_DEBUG_DATA_ROOT_DIR` or `/tmp/enos-debug-data`. |
| `setup-terraform` | no | `true` | When `true`, installs Terraform. |
| `setup-enos` | no | `true` | When `true`, installs Enos. |
| `launch-extra-args` | no | empty | Extra arguments added to the launch command before the scenario filter. |
| `retry-extra-args` | no | empty | Extra arguments added to the retry launch command before the scenario filter. |
| `destroy-extra-args` | no | empty | Extra arguments added to destroy before the scenario filter. |
| `pre-destroy-command` | no | empty | Optional enos scenario command to run before destroying (e.g., `exec`). |
| `pre-destroy-extra-args` | no | empty | Extra arguments added to the pre-destroy command before the scenario filter. |

## Outputs

| Name | Description |
| --- | --- |
| `launch-outcome` | Final main phase outcome: `success` or `failure`. |
| `launch-retry-attempted` | `true` when the main phase retry was attempted. |
| `destroy-outcome` | Final destroy outcome: `success` or `failure`. |
| `destroy-retry-attempted` | `true` when destroy retry was attempted. |
| `debug-data-artifact-name` | Sanitized debug artifact name derived from the scenario filter. |
| `debug-data-dir` | Debug data directory used by the action. |
| `error-message` | Failure message captured from the retry step, when present. |

## Usage

### Standard launch workflow

```yaml
- name: Run Enos scenario
  id: run-scenario
  uses: ./.github/actions/run-enos-scenario
  with:
    scenario-filter: ${{ matrix.scenario.id.filter }}
    timeout: 45m0s
    retry-timeout: 30m0s
    destroy-timeout: 10m0s
    show-state-on-retry: 'true'
```

### With custom destroy arguments

```yaml
- name: Run Enos scenario
  id: run-scenario
  uses: ./.github/actions/run-enos-scenario
  with:
    scenario-filter: ${{ inputs.scenario }}
    timeout: 60m0s
    retry-timeout: 60m0s
    destroy-timeout: 60m0s
    destroy-extra-args: --grpc-listen http://localhost
```

### Skip tool installation

```yaml
- name: Run Enos scenario
  uses: ./.github/actions/run-enos-scenario
  with:
    scenario-filter: ${{ steps.sample.outputs.filter }}
    setup-terraform: 'false'
    setup-enos: 'false'
```

### With pre-destroy command

```yaml
- name: Run Enos scenario
  uses: ./.github/actions/run-enos-scenario
  with:
    scenario-filter: ${{ inputs.scenario }}
    timeout: 30m0s
    pre-destroy-command: exec
    pre-destroy-extra-args: --cmd 'output -raw scan_markdown'
```

### With pre-destroy command and output redirection

```yaml
- name: Run Enos scenario
  uses: ./.github/actions/run-enos-scenario
  with:
    scenario-filter: ${{ inputs.scenario }}
    timeout: 30m0s
    pre-destroy-command: exec
    pre-destroy-extra-args: --cmd 'output -raw scan_markdown' | tee -a "$GITHUB_STEP_SUMMARY"
```
