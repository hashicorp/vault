/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { SELECTORS as PAGE } from 'vault/tests/helpers/kv/kv-page-selectors';

module('Integration | Component | kv | KvSecretForm', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.codeMirror = this.owner.lookup('service:code-mirror');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.store = this.owner.lookup('service:store');
    this.backend = 'my-kv-engine';
    this.path = 'my-secret';
    this.secret = this.store.createRecord('kv/data', { backend: this.backend });
    this.onSave = () => {};
    this.onCancel = () => {};
  });

  test('it makes post request on save', async function (assert) {
    assert.expect(2);
    this.server.post(`${this.backend}/data/${this.path}`, (schema, req) => {
      assert.ok(true, 'Request made to save secret');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, {
        data: {
          foo: 'bar',
        },
        options: {
          cas: 0,
        },
      });
    });
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');
    await render(
      hbs`
        <KvSecretForm
          @secret={{this.secret}}
          @onSave={{this.onSave}}
        />`,
      { owner: this.engine }
    );

    await fillIn(PAGE.form.inputByAttr('path'), this.path);
    await fillIn(PAGE.form.keyInput, 'foo');
    await fillIn(PAGE.form.valueInput, 'bar');
    await click(PAGE.form.secretSave);
  });
});
