/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const mockedEmptyResponse = {
  data: {
    auth_methods: {},
    kvv1_secrets: 0,
    kvv2_secrets: 0,
    lease_count_quotas: {},
    leases_by_auth_method: {},
    replication_status: {},
    secret_engines: {},
  },
};

export const mockedResponseWithData = {
  data: {
    auth_methods: { aws: 42, userpass: 43, kubernetes: 44 },
    kvv1_secrets: 60,
    kvv2_secrets: 40,
    lease_count_quotas: {
      global_lease_count_quota: { capacity: 420000, count: 210000, name: 'default' },
      total_lease_count_quotas: 1,
    },
    namespaces: 1,
    replication_status: {
      dr_primary: true,
      dr_state: 'enabled',
      pr_primary: false,
      pr_state: 'enabled',
    },
    secret_engines: { cubbyhole: 45, nomad: 46, aws: 47 },
  },
};
