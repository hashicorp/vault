module.exports = {
  env: {
    embertest: true,
  },
  globals: {
    faker: true,
    server: true,
    $: true,
    authLogout: false,
    authLogin: false,
    pollCluster: false,
    mountSupportedSecretBackend: false,
  },
};
