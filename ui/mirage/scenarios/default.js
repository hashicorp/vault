export default function(server) {
  server.createList('user', 10);
  server.create('user', { admin: true });
  server.create('metrics/activity');
  server.create('metrics/config');
}
