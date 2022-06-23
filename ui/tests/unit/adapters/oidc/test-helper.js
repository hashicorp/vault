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
    assert.expect(1);

    this.server.get(`/identity/${this.modelName}`, (schema, req) => {
      assert.equal(req.queryParams.list, 'true', 'request is made to correct endpoint on query');
      return { data: { keys: [this.data.name] } };
    });

    this.store.query(this.modelName, {});
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
