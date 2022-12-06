export default function (server) {
  server.get('keymgmt/key?list=true', function () {
    return {
      data: {
        keys: ['example-1', 'example-2', 'example-3'],
      },
    };
  });

  server.get('keymgmt/key/:name', function (_, request) {
    const name = request.params.name;
    return {
      data: {
        name,
        deletion_allowed: false,
        keys: {
          1: {
            creation_time: '2020-11-02T15:54:58.768473-08:00',
            public_key: '-----BEGIN PUBLIC KEY----- ... -----END PUBLIC KEY-----',
          },
          2: {
            creation_time: '2020-11-04T16:58:47.591718-08:00',
            public_key: '-----BEGIN PUBLIC KEY----- ... -----END PUBLIC KEY-----',
          },
        },
        latest_version: 2,
        min_enabled_version: 1,
        type: 'rsa-2048',
      },
    };
  });

  server.get('keymgmt/key/:name/kms', function () {
    return {
      data: {
        keys: ['example-kms'],
      },
    };
  });

  server.post('keymgmt/key/:name', function () {
    return {};
  });

  server.put('keymgmt/key/:name', function () {
    return {};
  });

  server.get('/keymgmt/kms/:provider/key', () => {
    const keys = [];
    let i = 1;
    while (i <= 75) {
      keys.push(`testkey-${i}`);
      i++;
    }
    return { data: { keys } };
  });
}
