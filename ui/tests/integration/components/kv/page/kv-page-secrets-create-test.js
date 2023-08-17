/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { hbs } from 'ember-cli-htmlbars';
import { click, fillIn, findAll, render, typeIn } from '@ember/test-helpers';
import codemirror from 'vault/tests/helpers/codemirror';
import { FORM } from 'vault/tests/helpers/kv/kv-selectors';
import sinon from 'sinon';

module('Integration | Component | kv-v2 | Page::Secrets::Create', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.router = this.owner.lookup('service:router');
    this.transitionStub = sinon.stub(this.router, 'transitionTo');
    this.backend = 'my-kv-engine';
    this.path = 'my-secret';
    this.maxVersions = 10;
    this.secret = this.store.createRecord('kv/data', { backend: this.backend, casVersion: 0 });
    this.metadata = this.store.createRecord('kv/metadata', { backend: this.backend });
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: 'create' },
    ];
  });

  hooks.afterEach(function () {
    this.router.transitionTo.restore();
  });

  test('it saves secret data and metadata', async function (assert) {
    assert.expect(5);
    this.server.post(`${this.backend}/data/${this.path}`, (schema, req) => {
      assert.ok(true, 'Request made to save secret');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, {
        data: { foo: 'bar' },
        options: { cas: 0 },
      });
      return {
        request_id: 'bd76db73-605d-fcbc-0dad-d44a008f9b95',
        data: {
          created_time: '2023-07-28T18:47:32.924809Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: false,
          version: 1,
        },
      };
    });

    this.server.post(`${this.backend}/metadata/${this.path}`, (schema, req) => {
      assert.ok(true, 'Request made to save metadata');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, {
        cas_required: false,
        custom_metadata: {
          'my-custom': 'metadata',
        },
        delete_version_after: '0s',
        max_versions: 10,
      });
    });

    await render(
      hbs`<Page::Secrets::Create
  @secret={{this.secret}}
  @metadata={{this.metadata}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    await fillIn(FORM.inputByAttr('path'), this.path);
    await fillIn(FORM.keyInput(), 'foo');
    await fillIn(FORM.maskedValueInput(), 'bar');

    await click(FORM.toggleMetadata);
    await fillIn(`[data-test-field="customMetadata"] ${FORM.keyInput()}`, 'my-custom');
    await fillIn(`[data-test-field="customMetadata"] ${FORM.valueInput()}`, 'metadata');
    await fillIn(FORM.inputByAttr('maxVersions'), this.maxVersions);

    await click(FORM.saveBtn);

    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.kv.secret.details'),
      'router transitions to secret details route on save'
    );
  });

  test('it does not send request to save secret metadata if fields are unchanged', async function (assert) {
    // this test contains two assertions, but only expects one because a request to kv/metadata
    // should NOT happen if its form inputs have not been edited
    assert.expect(1);
    this.server.post(`${this.backend}/data/${this.path}`, () => {
      assert.ok(true, 'Request only made to save secret');
      return {
        request_id: 'bd76db73-605d-fcbc-0dad-d44a008f9b95',
        data: {
          created_time: '2023-07-28T18:47:32.924809Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: false,
          version: 1,
        },
      };
    });

    this.server.post(`${this.backend}/metadata/${this.path}`, () => {
      // this assertion should not be hit!!
      assert.notOk(true, 'Request should not be made to save metadata');
      return new Response(403, {}, { errors: ['This request should not have been made'] });
    });

    await render(
      hbs`<Page::Secrets::Create
  @secret={{this.secret}}
  @metadata={{this.metadata}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    await fillIn(FORM.inputByAttr('path'), this.path);
    await fillIn(FORM.keyInput(), 'foo');
    await fillIn(FORM.maskedValueInput(), 'bar');
    await click(FORM.saveBtn);
  });

  test('it saves nested secrets', async function (assert) {
    assert.expect(3);
    const pathToSecret = 'path/to/secret/';
    this.secret.path = pathToSecret;
    this.server.post(`${this.backend}/data/${pathToSecret + this.path}`, (schema, req) => {
      assert.ok(true, 'Request made to save secret');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, {
        data: { foo: 'bar' },
        options: { cas: 0 },
      });
      return {
        request_id: 'bd76db73-605d-fcbc-0dad-d44a008f9b95',
        data: {
          created_time: '2023-07-28T18:47:32.924809Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: false,
          version: 1,
        },
      };
    });

    await render(
      hbs`<Page::Secrets::Create
  @secret={{this.secret}}
  @metadata={{this.metadata}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    assert.dom(FORM.inputByAttr('path')).hasValue(pathToSecret);
    await typeIn(FORM.inputByAttr('path'), this.path);
    await fillIn(FORM.keyInput(), 'foo');
    await fillIn(FORM.maskedValueInput(), 'bar');
    await click(FORM.saveBtn);
  });

  test('it renders API errors', async function (assert) {
    // this test contains an extra assertion because a request to kv/metadata
    // should NOT happen if kv/data fails
    assert.expect(3);
    this.server.post(`${this.backend}/data/${this.path}`, () => {
      return new Response(500, {}, { errors: ['nope'] });
    });

    this.server.post(`${this.backend}/metadata/${this.path}`, () => {
      // this assertion should not be hit because the request to save secret data failed!!
      assert.ok(true, 'Request made to save metadata');
      return new Response(403, {}, { errors: ['This request should not have been made'] });
    });

    await render(
      hbs`<Page::Secrets::Create
  @secret={{this.secret}}
  @metadata={{this.metadata}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    await fillIn(FORM.inputByAttr('path'), this.path);
    await click(FORM.saveBtn);
    assert.dom(FORM.messageError).hasText('Error nope', 'it renders API error');
    assert.dom(FORM.inlineAlert).hasText('There was an error submitting this form.');
    await click(FORM.cancelBtn);
    assert.ok(
      this.transitionStub.calledWith('vault.cluster.secrets.backend.kv.list'),
      'router transitions to secret list route on cancel'
    );
  });

  test('it renders kv secret validations', async function (assert) {
    assert.expect(6);

    await render(
      hbs`<Page::Secrets::Create
  @secret={{this.secret}}
  @metadata={{this.metadata}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    await typeIn(FORM.inputByAttr('path'), 'space ');
    assert
      .dom(FORM.validation('path'))
      .hasText(
        `Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.`
      );

    await fillIn(FORM.inputByAttr('path'), ''); // clear input
    await typeIn(FORM.inputByAttr('path'), 'slash/');
    assert.dom(FORM.validation('path')).hasText(`Path can't end in forward slash '/'.`);

    await typeIn(FORM.inputByAttr('path'), 'secret');
    assert
      .dom(FORM.validation('path'))
      .doesNotExist('it removes validation on key up when secret contains slash but does not end in one');

    await click(FORM.toggleJson);
    codemirror().setValue('i am a string and not JSON');
    assert
      .dom(FORM.inlineAlert)
      .hasText('JSON is unparsable. Fix linting errors to avoid data discrepancies.');

    codemirror().setValue('{}'); // clear linting error
    await fillIn(FORM.inputByAttr('path'), '');
    await click(FORM.saveBtn);
    const [pathValidation, formAlert] = findAll(FORM.inlineAlert);
    assert.dom(pathValidation).hasText(`Path can't be blank.`);
    assert.dom(formAlert).hasText('There is an error with this form.');
  });

  test('it toggles JSON view and saves modified data', async function (assert) {
    assert.expect(4);
    this.server.post(`${this.backend}/data/${this.path}`, (schema, req) => {
      assert.ok(true, 'Request made to save secret');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, {
        data: { hello: 'there' },
        options: { cas: 0 },
      });
      return {
        request_id: 'bd76db73-605d-fcbc-0dad-d44a008f9b95',
        data: {
          created_time: '2023-07-28T18:47:32.924809Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: false,
          version: 1,
        },
      };
    });

    await render(
      hbs`<Page::Secrets::Create
  @secret={{this.secret}}
  @metadata={{this.metadata}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    assert.dom(FORM.dataInputLabel({ isJson: false })).hasText('Secret data');
    await click(FORM.toggleJson);
    assert.dom(FORM.dataInputLabel({ isJson: true })).hasText('Secret data');

    codemirror().setValue(`{ "hello": "there"}`);
    await fillIn(FORM.inputByAttr('path'), this.path);
    await click(FORM.saveBtn);
  });
});
