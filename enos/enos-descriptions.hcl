// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

globals {
  description = {
    build_vault = <<-EOF
      Determine which Vault artifact we want to use for the scenario. Depending on the
      'artifact_source' variant we'll either build Vault from the local branch, fetch a candidate
      build from Artifactory, or use a local artifact that was built in CI via CRT.
    EOF

    create_backend_cluster = <<-EOF
      Create a storage backend cluster if necessary. When configured to use Consul it will
      install, configure, and start the Consul cluster on the target hosts and wait for the Consul
      cluster to become healthy. When using integrated raft storage this step is a no-op as the
      Vault cluster nodes will provide their own integrated storage.
    EOF

    create_seal_key = <<-EOF
      Create the necessary seal key infrastructure for Vaults auto-unseal functionality. Depending
      on the 'seal' variant this step will perform different actions. When using 'shamir' the step
      is a no-op as we won't require an external seal mechanism. When using 'pkcs11' this step will
      create a SoftHSM slot and associated token which can be distributed to all target nodes. When
      using 'awskms' a new AWSKMS key will be created. The necessary security groups and policies
      for Vault target nodes to access it the AWSKMS key are handled in the target modules.
    EOF

    create_vault_cluster = <<-EOF
      Create the the Vault cluster. In this module we'll install, configure, start, initialize and
      unseal all the nodes in the Vault. After initialization it also enables various audit engines.
    EOF

    create_vault_cluster_backend_targets = <<-EOF
      Create the target machines that we'll install Consul onto when using Consul for storage. We
      also handle creating AWS instance profiles and security groups that allow for auto-discovery
      via the retry_join functionality in Consul. The security group firewall rules will
      automatically allow SSH access from the host external IP address of the machine executing
      Enos, in addition to all of the required ports for Consul to function and be accessible in the
      VPC.
    EOF

    create_vault_cluster_targets = <<-EOF
      Create the target machines that we'll install Vault onto. We also handle creating AWS instance
      profiles and security groups that allow for auto-discovery via the retry_join functionality in
      Consul. The security group firewall rules will automatically allow SSH access from the host
      external IP address of the machine executing Enos, in addition to all of the required ports
      for Vault to function and be accessible in the VPC.
    EOF

    create_vpc = <<-EOF
      Create an AWS VPC, internet gateway, default security group, and default subnet that allows
      egress traffic via the internet gateway.
    EOF

    ec2_info = <<-EOF
      Query various endpoints in AWS Ec2 to gather metadata we'll use later in our run when creating
      infrastructure for the Vault cluster. This metadata includes:
        - AMI IDs for different Linux distributions and platform architectures
        - Available Ec2 Regions
        - Availability Zones for our desired machine instance types
    EOF

    enable_multiseal = <<-EOF
      Configure the Vault cluster with 'enable_multiseal' and up to three auto-unseal methods
      via individual, prioritized 'seal' stanzas.
    EOF

    get_local_metadata = <<-EOF
      Performs several Vault quality verification that are dynamically modified based on the Vault
      binary version, commit SHA, build-date (commit SHA date), and edition metadata. When we're
      testing existing artifacts this expected metadata is passed in via Enos variables. When we're
      building a local by using the 'artifact_source:local' variant, this step executes and
      populates the expected metadata with that of our branch so that we don't have to update the
      Enos variables on each commit.
    EOF

    get_vault_cluster_ip_addresses = <<-EOF
      Map the public and private IP addresses of the Vault cluster nodes and segregate them by
      their leader status. This allows us to easily determine the public IP addresses of the leader
      and follower nodes.
    EOF

    read_backend_license = <<-EOF
      When using Consul Enterprise as a storage backend, ensure that a Consul Enterprise license is
      present on disk and read its contents so that we can utilize it when configuring the storage
      cluster. Must have the 'backend:consul' and 'consul_edition:ent' variants.
    EOF

    read_vault_license = <<-EOF
      When deploying Vault Enterprise, ensure a Vault Enterprise license is present on disk and
      read its contents so that we can utilize it when configuring the Vault Enterprise cluster.
      Must have the 'edition' variant to be set to any Enterprise edition.
    EOF

    shutdown_nodes = <<-EOF
      Shut down the nodes to ensure that they are no longer operating software as part of the
      cluster.
    EOF

    start_vault_agent = <<-EOF
      Create an agent approle in the auth engine, generate a Vault Agent configuration file, and
      start the Vault agent.
    EOF

    stop_vault = <<-EOF
      Stop the Vault cluster by stopping the vault service via systemctl.
    EOF

    vault_leader_step_down = <<-EOF
      Force the Vault cluster leader to step down which forces the Vault cluster to perform a leader
      election.
    EOF

    verify_agent_output = <<-EOF
      Vault running in Agent mode uses templates to create log output.
    EOF

    verify_log_secrets = <<-EOF
      Verify that the vault audit log and systemd journal do not leak secret values.
    EOF

    verify_raft_cluster_all_nodes_are_voters = <<-EOF
      When configured with a 'backend:raft' variant, verify that all nodes in the cluster are
      healthy and are voters.
    EOF

    verify_autopilot_idle_state = <<-EOF
      Wait for the Autopilot to upgrade the entire Vault cluster and ensure that the target version
      matches the candidate version. Ensure that the cluster reaches an upgrade state of
      'await-server-removal'.
    EOF

    verify_replication_status = <<-EOF
      Verify that the default replication status is correct depending on the edition of Vault that
      been deployed. When testing a Community Edition of Vault we'll ensure that replication is not
      enabled. When testing any Enterprise edition of Vault we'll ensure that Performance and
      Disaster Recovery replication are available.
    EOF

    verify_seal_rewrap_entries_processed_eq_entries_succeeded_post_rewrap = <<-EOF
      Verify that the v1/sys/sealwrap/rewrap Vault API returns the rewrap data and
      'entries.processed' equals 'entries.succeeded' after the rewrap has completed.
    EOF

    verify_seal_rewrap_entries_processed_is_gt_zero_post_rewrap = <<-EOF
      Verify that the /sys/sealwrap/rewrap Vault API returns the rewrap data and the 'entries.processed' has
      processed at least one entry after the rewrap has completed.
    EOF

    verify_seal_rewrap_is_running_false_post_rewrap = <<-EOF
      Verify that the v1/sys/sealwrap/rewrap Vault API returns the rewrap data and 'is_running' is set to
      'false' after a rewrap has completed.
    EOF

    verify_seal_rewrap_no_entries_fail_during_rewrap = <<-EOF
      Verify that the v1/sys/sealwrap/rewrap Vault API returns the rewrap data and 'entries.failed' is '0'
      after the rewrap has completed.
    EOF

    verify_seal_type = <<-EOF
      Vault's reported seal type matches our configuration.
    EOF

    verify_secrets_engines_create = <<-EOF
      Verify that Vault is capable mounting, configuring, and using various secrets engines and auth
      methods. These currently include:
        - v1/auth/userpass/*
        - v1/identity/*
        - v1/kv/*
        - v1/sys/policy/*
    EOF

    verify_secrets_engines_read = <<-EOF
      Verify that data that we've created previously is still valid, consistent, and duarable.
      This includes:
        - v1/auth/userpass/*
        - v1/identity/*
        - v1/kv/*
        - v1/sys/policy/*
    EOF

    verify_ui = <<-EOF
      The Vault UI assets are embedded in the Vault binary and available when running.
    EOF

    verify_vault_unsealed = <<-EOF
      Verify that the Vault cluster has successfully unsealed.
    EOF

    verify_vault_version = <<-EOF
      Verify that the Vault CLI has the correct embedded version metadata and that the Vault Cluster
      verision history includes our expected version. The CLI metadata that is validated includes
      the Vault version, edition, build date, and any special prerelease metadata.
    EOF

    wait_for_cluster_to_have_leader = <<-EOF
      Wait for a leader election to occur before we proceed with any further quality verification.
    EOF

    wait_for_seal_rewrap = <<-EOF
      Wait for the Vault cluster seal rewrap process to complete.
    EOF

    verify_billing_start_date = <<-EOF
      Verify that the billing start date has successfully rolled over to the latest billing year if needed.
    EOF

  }
}
