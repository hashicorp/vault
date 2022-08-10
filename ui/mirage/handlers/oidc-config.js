export default function (server) {
  server.get('/identity/oidc/key/default', () => {
    return {
      request_id: '5eeb2a46-8726-7336-76ee-8daea2860624',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: {
        algorithm: 'RS256',
        allowed_client_ids: ['*'],
        rotation_period: 86400,
        verification_ttl: 86400,
      },
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  });

  server.get('/identity/oidc/assignment/allow_all', () => {
    return {
      request_id: '9a948dca-96bf-24ed-39f9-e2cdd94bf90f',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: {
        entity_ids: ['*'],
        group_ids: ['*'],
      },
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  });

  server.get('/identity/oidc/key', () => {
    return {
      request_id: 'key-list-id',
      data: {
        keys: ['default'],
      },
    };
  });

  server.get('/identity/entity/id', () => ({
    data: {
      key_info: { '1234-12345': { name: 'test-entity' } },
      keys: ['1234-12345'],
    },
  }));
  server.get('/identity/group/id', () => ({
    data: {
      key_info: { 'abcdef-123': { name: 'test-group' } },
      keys: ['abcdef-123'],
    },
  }));
}
