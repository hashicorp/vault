import { Response } from 'miragejs';

export default function (server) {
  const methods = ['totp', 'duo', 'okta', 'pingid'];
  const required = {
    totp: ['issuer'],
    duo: ['secret_key', 'integration_key', 'api_hostname'],
    okta: ['org_name', 'api_token'],
    pingid: ['settings_file_base64'],
  };

  const validate = (type, data, cb) => {
    if (!methods.includes(type)) {
      return new Response(400, {}, { errors: [`Method ${type} not found`] });
    }
    if (data) {
      const missing = required[type].reduce((params, key) => {
        if (!data[key]) {
          params.push(key);
        }
        return params;
      }, []);
      if (missing.length) {
        return new Response(400, {}, { errors: [`Missing required parameters: [${missing.join(', ')}]`] });
      }
    }
    return cb();
  };

  const dbKeyFromType = (type) => `mfa${type.charAt(0).toUpperCase()}${type.slice(1)}Methods`;

  const generateListResponse = (schema, isMethod) => {
    let records = [];
    if (isMethod) {
      methods.forEach((method) => {
        records.addObjects(schema.db[dbKeyFromType(method)].where({}));
      });
    } else {
      records = schema.db.mfaLoginEnforcements.where({});
    }
    // seed the db with a few records if none exist
    if (!records.length) {
      if (isMethod) {
        methods.forEach((type) => {
          records.push(server.create(`mfa-${type}-method`));
        });
      } else {
        records = server.createList('mfa-login-enforcement', 4).toArray();
      }
    }
    const dataKey = isMethod ? 'id' : 'name';
    const data = records.reduce(
      (resp, record) => {
        resp.key_info[record[dataKey]] = record;
        resp.keys.push(record[dataKey]);
        return resp;
      },
      {
        key_info: {},
        keys: [],
      }
    );
    return { data };
  };

  // list methods
  server.get('/identity/mfa/method/', (schema) => {
    return generateListResponse(schema, true);
  });
  // fetch method by id
  server.get('/identity/mfa/method/:id', (schema, { params: { id } }) => {
    let record;
    for (const method of methods) {
      record = schema.db[dbKeyFromType(method)].find(id);
      if (record) {
        break;
      }
    }
    // inconvenient when testing edit route to return a 404 on refresh since mirage memory is cleared
    // flip this variable to test 404 state if needed
    const shouldError = false;
    // create a new record so data is always returned
    if (!record && !shouldError) {
      return { data: server.create('mfa-totp-method') };
    }
    return !record ? new Response(404, {}, { errors: [] }) : { data: record };
  });
  // create method
  server.post('/identity/mfa/method/:type', (schema, { params: { type }, requestBody }) => {
    const data = JSON.parse(requestBody);
    return validate(type, data, () => {
      const record = server.create(`mfa-${type}-method`, data);
      return { data: { method_id: record.id } };
    });
  });
  // update method
  server.put('/identity/mfa/method/:type/:id', (schema, { params: { type, id }, requestBody }) => {
    const data = JSON.parse(requestBody);
    return validate(type, data, () => {
      schema.db[dbKeyFromType(type)].update(id, data);
      return {};
    });
  });
  // delete method
  server.delete('/identity/mfa/method/:type/:id', (schema, { params: { type, id } }) => {
    return validate(type, null, () => {
      schema.db[dbKeyFromType(type)].remove(id);
      return {};
    });
  });
  // list enforcements
  server.get('/identity/mfa/login-enforcement', (schema) => {
    return generateListResponse(schema);
  });
  // fetch enforcement by name
  server.get('/identity/mfa/login-enforcement/:name', (schema, { params: { name } }) => {
    const record = schema.db.mfaLoginEnforcements.findBy({ name });
    // inconvenient when testing edit route to return a 404 on refresh since mirage memory is cleared
    // flip this variable to test 404 state if needed
    const shouldError = false;
    // create a new record so data is always returned
    if (!record && !shouldError) {
      return { data: server.create('mfa-login-enforcement', { name }) };
    }
    return !record ? new Response(404, {}, { errors: [] }) : { data: record };
  });
  // create/update enforcement
  server.post('/identity/mfa/login-enforcement/:name', (schema, { params: { name }, requestBody }) => {
    const data = JSON.parse(requestBody);
    // at least one method id is required
    if (!data.mfa_method_ids?.length) {
      return new Response(400, {}, { errors: ['missing method ids'] });
    }
    // at least one of the following targets is required
    const required = [
      'auth_method_accessors',
      'auth_method_types',
      'identity_group_ids',
      'identity_entity_ids',
    ];
    let hasRequired = false;
    for (let key of required) {
      if (data[key]?.length) {
        hasRequired = true;
        break;
      }
    }
    if (!hasRequired) {
      return new Response(
        400,
        {},
        {
          errors: [
            'One of auth_method_accessors, auth_method_types, identity_group_ids, identity_entity_ids must be specified',
          ],
        }
      );
    }
    if (schema.db.mfaLoginEnforcements.findBy({ name })) {
      schema.db.mfaLoginEnforcements.update({ name }, data);
    } else {
      schema.db.mfaLoginEnforcements.insert(data);
    }
    return { ...data, id: data.name };
  });
  // delete enforcement
  server.delete('/identity/mfa/login-enforcement/:name', (schema, { params: { name } }) => {
    schema.db.mfaLoginEnforcements.remove({ name });
    return {};
  });
  // endpoints for target selection
  server.get('/identity/group/id', () => ({
    data: {
      key_info: { '34db6b52-591e-bc22-8af0-4add5e167326': { name: 'test-group' } },
      keys: ['34db6b52-591e-bc22-8af0-4add5e167326'],
    },
  }));
  server.get('/identity/group/id/:id', () => ({
    data: {
      id: '34db6b52-591e-bc22-8af0-4add5e167326',
      name: 'test-group',
    },
  }));
  server.get('/identity/entity/id', () => ({
    data: {
      key_info: { 'f831667b-7392-7a1c-c0fc-33d48cb1c57d': { name: 'test-entity' } },
      keys: ['f831667b-7392-7a1c-c0fc-33d48cb1c57d'],
    },
  }));
  server.get('/identity/entity/id/:id', () => ({
    data: {
      id: 'f831667b-7392-7a1c-c0fc-33d48cb1c57d',
      name: 'test-entity',
    },
  }));
  server.get('/sys/auth', () => ({
    data: {
      'userpass/': { accessor: 'auth_userpass_bb95c2b1', type: 'userpass' },
    },
  }));
}
