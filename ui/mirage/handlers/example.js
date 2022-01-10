// this can be removed at any time
// example to use when creating new mirage dev handlers
export default function (server) {
  server.get('/sys/namespaces', () => ({
    data: {
      keys: ['foo/', 'bar/'],
    },
  }));
}
