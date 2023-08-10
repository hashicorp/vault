/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export default (test) => {
  test('it should make request to correct endpoint on save', async function (assert) {
    assert.expect(1);

    this.server.post(this.path, () => {
      assert.ok(true, 'request made to correct endpoint on save');
    });

    const model = this.store.createRecord(this.modelName, this.data);
    await model.save();
  });

  test('it should throw error if attempting to createRecord with an existing name', async function (assert) {
    const { modelName, data } = this;
    this.store.pushPayload(modelName, { modelName, name: data.name });

    const model = this.store.createRecord(modelName, data);
    assert.rejects(model.save(), `Error: A record already exists with the name: ${data.name}`);
  });

  test('it should make request to correct endpoint on find', async function (assert) {
    assert.expect(1);

    this.server.get(this.path, () => {
      assert.ok(true, 'request is made to correct endpoint on find');
      return { data: this.data };
    });

    this.store.findRecord(this.modelName, this.data.name);
  });

  test('it should make request to correct endpoint on query', async function (assert) {
    const keyInfoModels = ['client', 'provider']; // these models have key_info on the LIST response
    const { name, ...otherAttrs } = this.data; // excludes name from key_info data
    const key_info = { [name]: { ...otherAttrs } };

    this.server.get(`/identity/${this.modelName}`, (schema, req) => {
      assert.strictEqual(req.queryParams.list, 'true', 'request is made to correct endpoint on query');
      if (keyInfoModels.some((model) => this.modelName.includes(model))) {
        return { data: { keys: [name], key_info } };
      } else {
        return { data: { keys: [name] } };
      }
    });

    this.store.query(this.modelName, {});
  });

  test('it should filter query when passed filterFor and paramKey', async function (assert) {
    const keyInfoModels = ['client', 'provider']; // these models have key_info on the LIST response
    const keys = ['model-1', 'model-2', 'model-3'];
    const key_info = {
      'model-1': {
        model_id: 'a123',
        key: 'test-key',
        access_token_ttl: '30m',
        id_token_ttl: '1h',
      },
      'model-2': {
        model_id: 'b123',
        key: 'test-key',
        access_token_ttl: '30m',
        id_token_ttl: '1h',
      },
      'model-3': {
        model_id: 'c123',
        key: 'test-key',
        access_token_ttl: '30m',
        id_token_ttl: '1h',
      },
    };

    this.server.get(`/identity/${this.modelName}`, () => {
      if (keyInfoModels.some((model) => this.modelName.includes(model))) {
        return { data: { keys, key_info } };
      } else {
        return { data: { keys: [this.data.name] } };
      }
    });

    // test passing 'paramKey' and 'filterFor' to query and filterListResponse in adapters/named-path.js works as expected
    if (keyInfoModels.some((model) => this.modelName.includes(model))) {
      let testQuery = ['*', 'a123'];
      await this.store
        .query(this.modelName, { paramKey: 'model_id', filterFor: testQuery })
        .then((resp) => assert.strictEqual(resp.length, 3, 'returns all models when ids include glob (*)'));

      testQuery = ['*'];
      await this.store
        .query(this.modelName, { paramKey: 'model_id', filterFor: testQuery })
        .then((resp) => assert.strictEqual(resp.length, 3, 'returns all models when glob (*) is only id'));

      testQuery = ['b123'];
      await this.store.query(this.modelName, { paramKey: 'model_id', filterFor: testQuery }).then((resp) => {
        assert.strictEqual(resp.length, 1, 'filters response and returns only matching id');

        assert.strictEqual(resp[0].name, 'model-2', 'response contains correct model');
      });

      testQuery = ['b123', 'c123'];
      await this.store.query(this.modelName, { paramKey: 'model_id', filterFor: testQuery }).then((resp) => {
        assert.strictEqual(resp.length, 2, 'filters response when passed multiple ids');
        resp.forEach((m) =>
          assert.ok(['model-2', 'model-3'].includes(m.id), `it filters correctly and included: ${m.id}`)
        );
      });

      await this.store
        .query(this.modelName, { paramKey: 'nonexistent_key', filterFor: testQuery })
        .then((resp) => assert.ok(resp.isLoaded, 'does not error when paramKey does not exist'));

      assert.rejects(
        this.store.query(this.modelName, { paramKey: 'model_id', filterFor: 'some-string' }),
        'throws assertion when filterFor is not an array'
      );
    } else {
      const testQuery = ['b123', 'c123'];
      await this.store
        .query(this.modelName, { paramKey: 'model_id', filterFor: testQuery })
        .then((resp) => assert.ok(resp.isLoaded, 'does not error when key_info does not exist'));
    }
  });

  test('it passes allowed_client_id only when the param exists', async function (assert) {
    const keyInfoModels = ['client', 'provider']; // these models have key_info on the LIST response
    const { name, ...otherAttrs } = this.data; // excludes name from key_info data
    const key_info = { [name]: { ...otherAttrs } };

    this.server.get(`/identity/${this.modelName}`, (schema, req) => {
      if (this.modelName === 'oidc/provider') {
        assert.propEqual(
          req.queryParams,
          { list: 'true', allowed_client_id: 'a123' },
          'request has allowed_client_id as query param'
        );
      } else {
        assert.propEqual(req.queryParams, { list: 'true' }, 'request only has `list` param');
      }
      if (keyInfoModels.some((model) => this.modelName.includes(model))) {
        return { data: { keys: [name], key_info } };
      } else {
        return { data: { keys: [name] } };
      }
    });

    // only /provider accepts an allowed_client_id
    if (this.modelName === 'oidc/provider') {
      this.store.query(this.modelName, { allowed_client_id: 'a123' });
    } else {
      this.store.query(this.modelName, {});
    }
  });

  test('it should make request to correct endpoint on delete', async function (assert) {
    assert.expect(1);

    this.server.get(this.path, () => ({ data: this.data }));
    this.server.delete(this.path, () => {
      assert.ok(true, 'request made to correct endpoint on delete');
    });

    const model = await this.store.findRecord(this.modelName, this.data.name);
    await model.destroyRecord();
  });
};
