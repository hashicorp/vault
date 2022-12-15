export default function (server, shouldConfigureRoles = true) {
  server.create('kubernetes-config', { path: 'kubernetes' });
  if (shouldConfigureRoles) {
    server.create('kubernetes-role');
    server.create('kubernetes-role', 'withRoleName');
    server.create('kubernetes-role', 'withRoleRules');
  }
}
