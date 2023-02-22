import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  kubernetes_host: 'https://192.168.99.100:8443',
  kubernetes_ca_cert:
    '-----BEGIN CERTIFICATE-----\nMIIDNTCCAh2gApGgAwIBAgIULNEk+01LpkDeJujfsAgIULNEkAgIULNEckApGgAwIBAg+01LpkDeJuj\n-----END CERTIFICATE-----',
  disable_local_ca_jwt: true,

  // property used only for record lookup and filtered from response payload
  path: null,
});
