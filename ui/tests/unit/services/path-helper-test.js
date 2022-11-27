import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

const openapiStub = {
  openapi: {
    components: {
      schemas: {
        UsersRequest: {
          type: 'object',
          properties: {
            password: {
              description: 'Password for the user',
              type: 'string',
              'x-vault-displayAttrs': { sensitive: true },
            },
          },
        },
      },
    },
    paths: {
      '/users/{username}': {
        post: {
          requestBody: {
            content: {
              'application/json': {
                schema: { $ref: '#/components/schemas/UsersRequest' },
              },
            },
          },
        },
        parameters: [
          {
            description: 'Username for this user.',
            in: 'path',
            name: 'username',
            required: true,
            schema: { type: 'string' },
          },
        ],
        'x-vault-displayAttrs': { itemType: 'User', action: 'Create' },
      },
    },
  },
};

module('Unit | Service | path-help', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.pathHelp = this.owner.lookup('service:path-help');
    this.store = this.owner.lookup('service:store');
  });

  test('it should generate model with mutableId', async function (assert) {
    assert.expect(2);

    this.server.get('/auth/userpass/', () => openapiStub);
    this.server.get('/auth/userpass/users/example', () => openapiStub);
    this.server.post('/auth/userpass/users/test', () => {
      assert.ok(true, 'POST request made to correct endpoint');
      return;
    });

    const modelType = 'generated-user-userpass';
    await this.pathHelp.getNewModel(modelType, 'userpass', 'auth/userpass/', 'user');
    const model = this.store.createRecord(modelType);
    model.set('mutableId', 'test');
    await model.save();
    assert.strictEqual(model.get('id'), 'test', 'model id is set to mutableId value on save success');
  });
});
