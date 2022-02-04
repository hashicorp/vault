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
