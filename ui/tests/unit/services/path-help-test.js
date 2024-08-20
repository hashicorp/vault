/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import Sinon from 'sinon';
import { reject } from 'rsvp';

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

  module('getNewModel', function (hooks) {
    hooks.beforeEach(function () {
      this.server.get('/auth/userpass/', () => openapiStub);
      this.server.get('/auth/userpass/users/example', () => openapiStub);
    });
    test('it generates a model with mutableId', async function (assert) {
      assert.expect(2);
      this.server.post('/auth/userpass/users/test', () => {
        assert.true(true, 'POST request made to correct endpoint');
        return;
      });

      const modelType = 'generated-user-userpass';
      await this.pathHelp.getNewModel(modelType, 'userpass', 'auth/userpass/', 'user');
      const model = this.store.createRecord(modelType);
      model.set('mutableId', 'test');
      await model.save();
      assert.strictEqual(model.get('id'), 'test', 'model id is set to mutableId value on save success');
    });

    test('it only generates the model once', async function (assert) {
      assert.expect(2);
      Sinon.spy(this.pathHelp, 'getPaths');

      const modelType = 'generated-user-userpass';
      await this.pathHelp.getNewModel(modelType, 'userpass', 'auth/userpass/', 'user');
      assert.true(this.pathHelp.getPaths.calledOnce, 'getPaths is called for new generated model');

      await this.pathHelp.getNewModel(modelType, 'userpass2', 'auth/userpass/', 'user');
      assert.true(this.pathHelp.getPaths.calledOnce, 'not called again even with different backend path');
    });

    test('it resolves without error if model already exists', async function (assert) {
      Sinon.stub(this.pathHelp, 'getPaths').callsFake(() => {
        assert.notOk(true, 'this method should not be called');
        return reject();
      });
      const modelType = 'kv/data';
      await this.pathHelp.getNewModel(modelType, 'my-kv').then(() => {
        assert.true(true, 'getNewModel resolves');
      });
    });
  });

  module('hydrateModel', function () {
    test('it should hydrate an existing model', async function (assert) {
      this.server.get(`/pki2/roles/example`, () => openapiStub);

      const modelType = 'pki/role';
      await this.pathHelp.hydrateModel(modelType, 'pki2');
      const model = this.store.createRecord(modelType);
      model.set('username', 'foobar');
      assert.strictEqual(model.username, 'foobar', 'sets value of key that only exists in openAPI response');
    });
  });
});
