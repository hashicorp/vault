/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { click, render, waitFor } from '@ember/test-helpers';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import KvForm from 'vault/forms/secrets/kv';

module('Integration | Component | kv-v2 | Page::Secrets::Create', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(function () {
    this.backend = 'my-kv-engine';
    this.path = 'my-secret';
    this.form = new KvForm(
      {
        path: this.path,
        max_versions: 0,
        delete_version_after: '0s',
        cas_required: false,
      },
      { isNew: true }
    );
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: 'Create' },
    ];

    this.renderComponent = () =>
      render(
        hbs`
        <Page::Secrets::Create
          @form={{this.form}}
          @path={{this.path}}
          @backend={{this.backend}}
          @breadcrumbs={{this.breadcrumbs}}
        />
      `,
        { owner: this.engine }
      );

    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    setRunOptions({
      rules: {
        // TODO fix JSONEditor, KVObjectEditor, MaskedInput
        label: { enabled: false },
        'color-contrast': { enabled: false }, // JSONEditor only
      },
    });
  });

  test('it should toggle json editor', async function (assert) {
    assert.expect(4);

    await this.renderComponent();

    assert.dom(FORM.dataInputLabel({ isJson: false })).hasText('Secret data');
    assert.dom('.cm-editor').doesNotExist('CodeMirror editor is not rendered');

    await click(GENERAL.toggleInput('json'));
    assert.dom(FORM.dataInputLabel({ isJson: true })).hasText('Secret data');
    await waitFor('.cm-editor');
    assert.dom('.cm-editor').exists('CodeMirror editor is rendered');
  });

  test('it should toggle metadata', async function (assert) {
    await this.renderComponent();

    assert.dom(FORM.toggleMetadata).hasText('Show secret metadata');
    assert.dom(PAGE.create.metadataSection).doesNotExist('metadata section is hidden');
    await click(FORM.toggleMetadata);
    assert.dom(FORM.toggleMetadata).hasText('Hide secret metadata');
    assert.dom(PAGE.create.metadataSection).exists('metadata section is shown');
  });
});
