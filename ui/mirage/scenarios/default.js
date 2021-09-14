export default function(server) {
  server.create('metrics/config');
  server.create('feature', { feature_flags: ['SOME_FLAG', 'VAULT_CLOUD_ADMIN_NAMESPACE'] });
}
