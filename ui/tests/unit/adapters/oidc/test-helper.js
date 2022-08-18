export default (test) => {
  test('it should make request to correct endpoint on save', async function (assert) {
    assert.expect(2);

    this.server.post(this.path, () => {
      assert.ok(true, 'request made to correct endpoint on save');
    });

    const model = this.store.createRecord(this.modelName, this.data);
    await model.save();
    await model.save();
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
    // these models have key_info on the LIST response
    const keyInfoModels = ['client', 'provider'];
    if (keyInfoModels.some((model) => this.modelName.includes(model))) {
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

      this.server.get(`/identity/${this.modelName}`, (schema, req) => {
        assert.equal(req.queryParams.list, 'true', 'request is made to correct endpoint on query');
        return { data: { keys, key_info } };
      });

      // test passing 'paramKey' and 'filterFor' to query and filterListResponse in adapters/named-path.js works as expected
      let testQuery = ['*', 'a123'];
      await this.store
        .query(this.modelName, { paramKey: 'model_id', filterFor: testQuery })
        .then((resp) =>
          assert.equal(resp.content.length, 3, 'returns all clients when ids include glob (*)')
        );

      testQuery = ['*'];
      await this.store
        .query(this.modelName, { paramKey: 'model_id', filterFor: testQuery })
        .then((resp) => assert.equal(resp.content.length, 3, 'returns all clients when glob (*) is only id'));

      testQuery = ['b123'];
      await this.store.query(this.modelName, { paramKey: 'model_id', filterFor: testQuery }).then((resp) => {
        console.log(resp);
        assert.equal(resp.content.length, 1, 'filters response and returns only matching id');
        console.log(resp.firstObject, 'FIRST OBJECT');
        assert.equal(resp.firstObject.name, 'model-2', 'response contains correct model');
      });

      testQuery = ['b123', 'c123'];
      await this.store.query(this.modelName, { paramKey: 'model_id', filterFor: testQuery }).then((resp) => {
        assert.equal(resp.content.length, 2, 'filters response when passed multiple ids');
        resp.content.forEach((m) =>
          assert.ok(['model-2', 'model-3'].includes(m.id), `it filters correctly and included: ${m.id}`)
        );
      });
    } else {
      this.server.get(`/identity/${this.modelName}`, (schema, req) => {
        assert.equal(req.queryParams.list, 'true', 'request is made to correct endpoint on query');
        return { data: { keys: [this.data.name] } };
      });

      this.store.query(this.modelName, {});

      this.store
        .query(this.modelName, { paramKey: 'model_id', filterFor: this.data.name })
        .then((resp) => assert.ok(resp.isLoaded, 'does not attempt to filter when no key_info'));
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
