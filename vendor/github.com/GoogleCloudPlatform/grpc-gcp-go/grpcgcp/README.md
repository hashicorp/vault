
## How to test Spanner integration

1. Set GCP project id with GCP_PROJECT_ID environment variable.

        export GCP_PROJECT_ID=test-project

1. Set service key credentials file using GOOGLE_APPLICATION_CREDENTIALS env variable.

        export GOOGLE_APPLICATION_CREDENTIALS=/service/account/credentials.json

1. Run the tests.

        go test -v

To skip Spanner setup run

        SKIP_SPANNER=true go test -v
