/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { hbs } from 'ember-cli-htmlbars';
import { click, fillIn, render } from '@ember/test-helpers';
import codemirror from 'vault/tests/helpers/codemirror';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | kv-v2 | Page::Secret::Edit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.router = this.owner.lookup('service:router');
    this.transitionStub = sinon.stub(this.router, 'transitionTo');
    this.backend = 'my-kv-engine';
    this.path = 'my-secret';
    this.secret = this.store.createRecord('kv/data', {
      backend: this.backend,
      path: this.path,
      secretData: { foo: 'bar' },
      casVersion: 1,
    });
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: 'Edit' },
    ];
    setRunOptions({
      rules: {
        // TODO fix JSONEditor, KVObjectEditor, MaskedInput
        label: { enabled: false },
        'color-contrast': { enabled: false }, // JSONEditor only
      },
    });
  });

  hooks.afterEach(function () {
    this.router.transitionTo.restore();
  });

  test('it saves a new secret version', async function (assert) {
    assert.expect(10);
    this.server.post(`${this.backend}/data/${this.path}`, (schema, req) => {
      assert.true(true, 'Request made to save secret');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(
        payload,
        { data: { foo: 'bar', foo2: 'bar2' }, options: { cas: 1 } },
        'request has expected payload'
      );
      return {
        request_id: 'bd76db73-605d-fcbc-0dad-d44a008f9b95',
        data: {
          created_time: '2023-07-28T18:47:32.924809Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: false,
          version: 2,
        },
      };
    });

    await render(
      hbs`<Page::Secret::Edit
  @secret={{this.secret}}
  @previousVersion={{4}}
  @currentVersion={{4}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    assert.dom(FORM.inputByAttr('path')).isDisabled();
    assert.dom(FORM.inputByAttr('path')).hasValue(this.path);
    assert.dom(FORM.keyInput()).hasValue('foo');
    assert.dom(FORM.maskedValueInput()).hasValue('bar');
    assert.dom(FORM.dataInputLabel({ isJson: false })).hasText('Version data');
    await click(FORM.toggleJson);
    assert.strictEqual(
      codemirror().getValue(' '),
      `{   \"foo": \"bar" }`, // eslint-disable-line no-useless-escape
      'json editor initializes with empty object'
    );
    assert.dom(FORM.dataInputLabel({ isJson: true })).hasText('Version data');
    await click(FORM.toggleJson);
    await fillIn(FORM.keyInput(1), 'foo2');
    await fillIn(FORM.maskedValueInput(1), 'bar2');
    await click(FORM.saveBtn);
    const [actual] = this.transitionStub.lastCall.args;
    assert.strictEqual(
      actual,
      'vault.cluster.secrets.backend.kv.secret.index',
      'router transitions to secret overview route on save'
    );
  });

  test('diff works correctly', async function (assert) {
    await render(
      hbs`<Page::Secret::Edit
  @secret={{this.secret}}
  @previousVersion={{4}}
  @currentVersion={{4}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    assert.dom(PAGE.edit.toggleDiff).isNotDisabled('Diff toggle is not disabled');
    assert.dom(PAGE.edit.toggleDiffDescription).hasText('No changes to show. Update secret to view diff');
    assert.dom(PAGE.diff.visualDiff).doesNotExist('Does not show visual diff');

    await fillIn(FORM.keyInput(1), 'foo2');
    await fillIn(FORM.maskedValueInput(1), 'bar2');

    assert.dom(PAGE.edit.toggleDiff).isNotDisabled('Diff toggle is not disabled');
    assert.dom(PAGE.edit.toggleDiffDescription).hasText('Showing the diff will reveal secret values');
    assert.dom(PAGE.diff.visualDiff).doesNotExist('Does not show visual diff');
    await click(PAGE.edit.toggleDiff);
    assert.dom(PAGE.diff.visualDiff).exists('Shows visual diff');
    assert.dom(PAGE.diff.added).hasText(`foo2"bar2"`);

    await click(FORM.toggleJson);
    codemirror().setValue('{ "foo3": "bar3" }');

    assert.dom(PAGE.diff.visualDiff).exists('Visual diff updates');
    assert.dom(PAGE.diff.deleted).hasText(`foo"bar"`);
    assert.dom(PAGE.diff.added).hasText(`foo3"bar3"`);
  });

  test('it saves nested secrets', async function (assert) {
    assert.expect(3);
    const nestedSecret = 'path/to/secret';
    this.secret.path = nestedSecret;
    this.server.post(`${this.backend}/data/${nestedSecret}`, (schema, req) => {
      assert.ok(true, 'Request made to save secret');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, {
        data: { foo: 'bar' },
        options: { cas: 1 },
      });
      return {
        request_id: 'bd76db73-605d-fcbc-0dad-d44a008f9b95',
        data: {
          created_time: '2023-07-28T18:47:32.924809Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: false,
          version: 2,
        },
      };
    });

    await render(
      hbs`<Page::Secret::Edit
  @secret={{this.secret}}
  @previousVersion={{4}}
  @currentVersion={{4}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    assert.dom(FORM.inputByAttr('path')).hasValue(nestedSecret);
    await click(FORM.saveBtn);
  });

  test('it renders API errors', async function (assert) {
    assert.expect(3);
    this.server.post(`${this.backend}/data/${this.path}`, () => {
      return new Response(500, {}, { errors: ['nope'] });
    });

    await render(
      hbs`<Page::Secret::Edit
  @secret={{this.secret}}
  @previousVersion={{4}}
  @currentVersion={{4}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    await click(FORM.saveBtn);
    assert.dom(FORM.messageError).hasText('Error nope', 'it renders API error');
    assert.dom(FORM.inlineAlert).hasText('There was an error submitting this form.');
    await click(FORM.cancelBtn);
    const [actual] = this.transitionStub.lastCall.args;
    assert.strictEqual(
      actual,
      'vault.cluster.secrets.backend.kv.secret.index',
      'router transitions to secret overview route on cancel'
    );
  });

  test('it renders kv secret validations', async function (assert) {
    assert.expect(2);

    await render(
      hbs`<Page::Secret::Edit
  @secret={{this.secret}}
  @previousVersion={{4}}
  @currentVersion={{4}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    await click(FORM.toggleJson);
    codemirror().setValue('i am a string and not JSON');
    assert
      .dom(FORM.inlineAlert)
      .hasText('JSON is unparsable. Fix linting errors to avoid data discrepancies.');

    codemirror().setValue(`""`);
    await click(FORM.saveBtn);
    assert.dom(FORM.inlineAlert).hasText('Vault expects data to be formatted as an JSON object.');
  });

  test('it toggles JSON view and saves modified data', async function (assert) {
    assert.expect(4);

    this.server.post(`${this.backend}/data/${this.path}`, (schema, req) => {
      assert.ok(true, 'Request made to save secret');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, {
        data: { hello: 'there' },
        options: { cas: 1 },
      });
      return {
        request_id: 'bd76db73-605d-fcbc-0dad-d44a008f9b95',
        data: {
          created_time: '2023-07-28T18:47:32.924809Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: false,
          version: 2,
        },
      };
    });

    await render(
      hbs`<Page::Secret::Edit
  @secret={{this.secret}}
  @previousVersion={{3}}
  @currentVersion={{4}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );
    assert.dom(FORM.dataInputLabel({ isJson: false })).hasText('Version data');
    await click(FORM.toggleJson);
    assert.dom(FORM.dataInputLabel({ isJson: true })).hasText('Version data');

    codemirror().setValue(`{ "hello": "there"}`);
    await click(FORM.saveBtn);
  });

  test('it renders alert when creating a new secret version from an old version', async function (assert) {
    assert.expect(1);

    await render(
      hbs`<Page::Secret::Edit
  @secret={{this.secret}}
  @previousVersion={{1}}
  @currentVersion={{4}}
  @breadcrumbs={{this.breadcrumbs}}
/>`,
      { owner: this.engine }
    );

    assert
      .dom(FORM.versionAlert)
      .hasText(
        `Warning You are creating a new version based on data from Version 1. The current version for my-secret is Version 4.`
      );
  });
});
