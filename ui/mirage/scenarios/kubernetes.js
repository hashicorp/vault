export default function (server) {
  server.create('kubernetes-config', { path: 'kubernetes' });
  server.create('kubernetes-role');
  server.create('kubernetes-role', 'withRoleName');
  server.create('kubernetes-role', 'withRoleRules');
}
