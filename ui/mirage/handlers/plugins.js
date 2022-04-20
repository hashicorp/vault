export default function (server) {
  server.get('/plugins', function () {
    return {
      data: { keys: ['my-plugin-1', 'slack'] },
    };
  });

  server.get('/plugin/:pluginname', function (db, req) {
    const name = req.params.pluginname;
    return {
      data: {
        name,
        type: 'my-custom-plugin-type',
        pages: [
          {
            url: 'http://localhost:3000/teams/linkedin/recruiting',
            tabName: 'LinkedIn',
            description: 'Any content goes here, but it cannot be html',
          },
          {
            url: 'http://localhost:3000/teams/ms/tbt',
            tabName: 'Microsoft',
            description: 'This is where you can do things related to Microsoft',
          },
          {
            url: 'http://localhost:3000/example/overview',
            tabName: 'Session test',
          },
        ],
      },
    };
  });
}
