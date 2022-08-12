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
      const key_info = [
        {
          'model-1': {
            key: 'test-key',
            access_token_ttl: '30m',
            id_token_ttl: '1h',
          },
        },
        {
          'model-2': {
            key: 'test-key',
            access_token_ttl: '30m',
            id_token_ttl: '1h',
          },
        },
        {
          'model-3': {
            key: 'test-key',
            access_token_ttl: '30m',
            id_token_ttl: '1h',
          },
        },
      ];

      this.server.get(`/identity/${this.modelName}`, (schema, req) => {
        assert.equal(req.queryParams.list, 'true', 'request is made to correct endpoint on query');
        return { data: { keys, key_info } };
      });

      let testQuery = ['*', 'model-1'];
      await this.store
        .query(this.modelName, { filterIds: testQuery })
        .then((resp) =>
          assert.equal(resp.content.length, 3, 'returns all clients when ids include glob (*)')
        );

      testQuery = ['*'];
      await this.store
        .query(this.modelName, { filterIds: testQuery })
        .then((resp) => assert.equal(resp.content.length, 3, 'returns all clients when glob (*) is only id'));

      testQuery = ['model-2'];
      await this.store.query(this.modelName, { filterIds: testQuery }).then((resp) => {
        assert.equal(resp.content.length, 1, 'filters response and returns only matching id');
        assert.equal(resp.firstObject.name, 'model-2', 'response contains correct model');
      });

      testQuery = ['model-2', 'model-3'];
      await this.store.query(this.modelName, { filterIds: testQuery }).then((resp) => {
        assert.equal(resp.content.length, 2, 'filters response when passed multiple ids');
      });
    } else {
      this.server.get(`/identity/${this.modelName}`, (schema, req) => {
        assert.equal(req.queryParams.list, 'true', 'request is made to correct endpoint on query');
        return { data: { keys: [this.data.name] } };
      });

      this.store.query(this.modelName, {});

      this.store
        .query(this.modelName, { filterIds: this.data.name })
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
