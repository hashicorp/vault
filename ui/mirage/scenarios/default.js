import ENV from 'vault/config/environment';
const { handler } = ENV['ember-cli-mirage'];

export default function (server) {
  server.create('clients/config');
  server.create('feature', { feature_flags: ['SOME_FLAG', 'VAULT_CLOUD_ADMIN_NAMESPACE'] });

  if (handler === 'kubernetes') {
    server.create('kubernetes-config', { path: 'kubernetes' });
    server.createList('kubernetes-role', 5);
  }
}
