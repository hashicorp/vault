#!/usr/bin/env bash
set -euo pipefail

ARTIFACT_NAME=$1

# Create Enos scenario test matrix based on the artifact to test
if [[ $ARTIFACT_NAME == *"ent"* ]]; then
  export enos_matrix="{\"include\":[{\"scenario\":\"smoke backend:consul consul_version:1.12.3 distro:ubuntu seal:awskms arch:amd64 builder:crt edition:ent\",\"aws_region\":\"us-west-1\"},{\"scenario\":\"smoke backend:raft consul_version:1.12.3 distro:ubuntu seal:shamir arch:amd64 builder:crt edition:ent\",\"aws_region\":\"us-west-2\"},{\"scenario\":\"upgrade backend:raft consul_version:1.11.7 distro:rhel seal:shamir arch:amd64 builder:crt edition:ent\",\"aws_region\":\"us-west-1\"},{\"scenario\":\"upgrade backend:consul consul_version:1.11.7 distro:rhel seal:awskms arch:amd64 builder:crt edition:ent\",\"aws_region\":\"us-west-2\"},{\"scenario\":\"autopilot distro:ubuntu seal:shamir arch:amd64 builder:crt edition:ent\",\"aws_region\":\"us-west-1\"}]}"
else
  export enos_matrix="{\"include\":[{\"scenario\":\"smoke backend:consul consul_version:1.12.3 distro:ubuntu seal:awskms arch:amd64 builder:crt edition:oss\",\"aws_region\":\"us-west-1\"},{\"scenario\":\"smoke backend:raft consul_version:1.12.3 distro:ubuntu seal:shamir arch:amd64 builder:crt edition:oss\",\"aws_region\":\"us-west-2\"},{\"scenario\":\"upgrade backend:raft consul_version:1.11.7 distro:rhel seal:shamir arch:amd64 builder:crt edition:oss\",\"aws_region\":\"us-west-1\"},{\"scenario\":\"upgrade backend:consul consul_version:1.11.7 distro:rhel seal:awskms arch:amd64 builder:crt edition:oss\",\"aws_region\":\"us-west-2\"}]}"
fi

echo "::set-output name=matrix::$enos_matrix"
