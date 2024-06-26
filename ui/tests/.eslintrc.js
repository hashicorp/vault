/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable no-undef */
module.exports = {
  env: {
    embertest: true,
  },
  globals: {
    server: true,
    $: true,
    authLogout: false,
    authLogin: false,
    pollCluster: false,
    mountSupportedSecretBackend: false,
  },
};
