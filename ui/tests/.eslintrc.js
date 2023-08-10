/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
