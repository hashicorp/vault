/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import Sinon from 'sinon';
import { reject } from 'rsvp';
import { Response } from 'miragejs';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

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

  hooks.afterEach(function () {
    Sinon.restore();
  });

  test('schema properties preserve identifiers, sensitivity, groups and legacy field names', async function (assert) {
    this.server.get('/auth/schema/users/example', () => ({
      ...openapiStub,
      openapi: {
        ...openapiStub.openapi,
        components: {
          schemas: {
            UsersRequest: {
              type: 'object',
              properties: {
                ...openapiStub.openapi.components.schemas.UsersRequest.properties,
                token_policies: {
                  type: 'array',
                  items: { type: 'string' },
                  description: 'Policies applied to generated tokens.',
                  'x-vault-displayAttrs': { name: 'Token policies', group: 'Tokens' },
                },
                token_ttl: { type: 'integer', format: 'seconds' },
                obsolete_policy: { type: 'string', deprecated: true },
              },
            },
          },
        },
      },
    }));

    const props = await this.pathHelp.getProps('/v1/auth/schema/users/example?help=true');
    assert.deepEqual(
      Object.keys(props),
      ['username', 'password', 'tokenPolicies', 'tokenTtl'],
      'deprecated fields are omitted and identifiers come first'
    );
    assert.strictEqual(props.username.fieldValue, 'mutableId', 'path identifier uses the editable record ID');
    assert.true(props.username.readOnly, 'path identifier is omitted from the request body');
    assert.true(props.password.sensitive, 'password remains masked');
    assert.deepEqual(
      props.tokenPolicies,
      {
        editType: 'stringArray',
        helpText: 'Policies applied to generated tokens.',
        fieldGroup: 'Tokens',
        label: 'Token policies',
      },
      'array field metadata preserves its group and label'
    );
    assert.strictEqual(props.tokenTtl.editType, 'ttl', 'duration fields use TTL inputs');
    assert.false('type' in props.tokenTtl, 'duration values are not coerced by a number transform');
  });

  test('inline schemas produce the same fields as referenced schemas', async function (assert) {
    this.server.get('/auth/schema/reference', () => openapiStub);
    this.server.get('/auth/schema/inline', () => ({
      openapi: {
        paths: {
          '/users/{username}': {
            ...openapiStub.openapi.paths['/users/{username}'],
            post: {
              requestBody: {
                content: {
                  'application/json': { schema: openapiStub.openapi.components.schemas.UsersRequest },
                },
              },
            },
          },
        },
      },
    }));

    const referenced = await this.pathHelp.getProps('/v1/auth/schema/reference?help=true');
    const inline = await this.pathHelp.getProps('/v1/auth/schema/inline?help=true');
    assert.deepEqual(inline, referenced, 'both schema shapes preserve the same field contract');
  });

  test('schema fetch failures are propagated to the caller', async function (assert) {
    this.server.get('/auth/schema/denied', () => new Response(403, {}, { errors: ['permission denied'] }));

    await assert.rejects(
      this.pathHelp.getProps('/v1/auth/schema/denied?help=true'),
      /permission denied/,
      'failed schema requests do not return success-shaped metadata'
    );
  });

  test('form metadata preserves API field names, field options and ordered groups without model registration', async function (assert) {
    this.server.get('/auth/schema/forms', () => ({
      ...openapiStub,
      openapi: {
        ...openapiStub.openapi,
        components: {
          schemas: {
            UsersRequest: {
              type: 'object',
              properties: {
                token_policies: {
                  type: 'array',
                  items: { type: 'string' },
                  description: 'Policies applied to generated tokens.',
                  'x-vault-displayAttrs': { name: 'Token policies', group: 'Tokens' },
                },
                ...openapiStub.openapi.components.schemas.UsersRequest.properties,
                token_ttl: { type: 'integer', format: 'seconds' },
                max_uses: {
                  type: 'integer',
                  'x-vault-displayAttrs': { value: 5, group: 'Tokens' },
                },
                obsolete_policy: { type: 'string', deprecated: true },
              },
            },
          },
        },
      },
    }));
    Sinon.spy(this.pathHelp, '_registerModel');

    const { formFields, formFieldGroups } = await this.pathHelp.getFormMetadata('/v1/auth/schema/forms');
    assert.deepEqual(
      formFields.map(({ name }) => name),
      ['username', 'token_policies', 'password', 'token_ttl', 'max_uses'],
      'API casing and path-identifier ordering are preserved; deprecated fields are omitted'
    );
    assert.true(
      formFields.every((field) => field instanceof FormField),
      'returns form field instances'
    );
    assert.true(
      formFieldGroups.every((group) => group instanceof FormFieldGroup),
      'returns form group instances'
    );
    assert.deepEqual(
      formFieldGroups.map((group) => Object.keys(group)[0]),
      ['default', 'Tokens'],
      'default group stays first even when token properties precede default properties'
    );
    assert.deepEqual(
      formFieldGroups[1].Tokens.map(({ name }) => name),
      ['token_policies', 'max_uses'],
      'grouped fields preserve their schema order'
    );
    const [username, policies, password, ttl, maxUses] = formFields;
    assert.strictEqual(username.options.fieldValue, 'username', 'identifier binds to plain form data');
    assert.strictEqual(policies.options.fieldValue, 'token_policies', 'API property is the field binding');
    assert.strictEqual(policies.options.label, 'Token policies', 'display label is preserved');
    assert.strictEqual(policies.options.editType, 'stringArray', 'array input behavior is preserved');
    assert.strictEqual(
      policies.options.subText,
      'Policies applied to generated tokens.',
      'form descriptions use subText'
    );
    assert.true(password.options.sensitive, 'password remains masked');
    assert.strictEqual(ttl.options.editType, 'ttl', 'duration fields use TTL inputs');
    assert.strictEqual(ttl.type, undefined, 'duration values are not coerced by a number transform');
    assert.strictEqual(maxUses.type, 'number', 'integer properties use number fields');
    assert.strictEqual(maxUses.options.defaultValue, 5, 'schema defaults are available to form consumers');
    assert.true(this.pathHelp._registerModel.notCalled, 'metadata does not generate an Ember Data model');
  });

  test('inline and referenced schemas produce equivalent form metadata', async function (assert) {
    this.server.get('/auth/schema/reference', () => openapiStub);
    this.server.get('/auth/schema/inline', () => ({
      openapi: {
        paths: {
          '/users/{username}': {
            ...openapiStub.openapi.paths['/users/{username}'],
            post: {
              requestBody: {
                content: {
                  'application/json': { schema: openapiStub.openapi.components.schemas.UsersRequest },
                },
              },
            },
          },
        },
      },
    }));
    assert.deepEqual(
      await this.pathHelp.getFormMetadata('/v1/auth/schema/inline'),
      await this.pathHelp.getFormMetadata('/v1/auth/schema/reference'),
      'both schema shapes preserve the same form field and group contract'
    );
  });

  test('form metadata rejects denied requests and missing request schemas', async function (assert) {
    this.server.get('/auth/schema/denied', () => new Response(403, {}, { errors: ['permission denied'] }));
    this.server.get('/auth/schema/missing', () => ({ openapi: { paths: {} } }));
    this.server.get('/auth/schema/broken-reference', () => ({
      ...openapiStub,
      openapi: { ...openapiStub.openapi, components: { schemas: {} } },
    }));
    await assert.rejects(
      this.pathHelp.getFormMetadata('/v1/auth/schema/denied'),
      /permission denied/,
      'denied metadata requests propagate the backend error'
    );
    await assert.rejects(
      this.pathHelp.getFormMetadata('/v1/auth/schema/missing'),
      /OpenAPI request schema not found/,
      'missing schemas do not silently return an empty form'
    );
    await assert.rejects(
      this.pathHelp.getFormMetadata('/v1/auth/schema/broken-reference'),
      /OpenAPI request properties not found/,
      'unresolved schema references fail explicitly'
    );
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

    test('generated fields serialize complete payloads without path identifiers or untouched values', async function (assert) {
      this.server.post('/auth/userpass/users/payload-user', (schema, request) => {
        assert.deepEqual(
          JSON.parse(request.requestBody),
          { password: 'disposable-password' },
          'payload includes the edited password but not the path identifier'
        );
        return {};
      });

      await this.pathHelp.getNewModel('generated-user-userpass', 'userpass', 'auth/userpass/', 'user');
      const model = this.store.createRecord('generated-user-userpass', {
        password: 'disposable-password',
      });
      model.set('mutableId', 'payload-user');
      await model.save();
      assert.strictEqual(model.id, 'payload-user', 'successful create assigns the path identifier');
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
      const modelType = 'cluster';
      await this.pathHelp.getNewModel(modelType, 'my-kv').then(() => {
        assert.true(true, 'getNewModel resolves');
      });
    });
  });
});
