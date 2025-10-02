/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { click, fillIn, render, settled, waitFor } from '@ember/test-helpers';
import codemirror, { setCodeEditorValue } from 'vault/tests/helpers/codemirror';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import KvForm from 'vault/forms/secrets/kv';

module('Integration | Component | kv-v2 | Page::Secret::Edit', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'my-kv-engine';
    this.path = 'my-secret';
    this.secret = {
      secretData: { foo: 'bar' },
      version: 1,
    };
    this.metadata = { current_version: 1 };
    this.form = new KvForm({
      path: this.path,
      secretData: this.secret.secretData,
      max_versions: 0,
      options: {
        cas: this.secret.version,
      },
    });
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: 'Edit' },
    ];

    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'kvV2Write').resolves();

    this.renderComponent = () =>
      render(
        hbs`
          <Page::Secret::Edit
            @form={{this.form}}
            @secret={{this.secret}}
            @metadata={{this.metadata}}
            @path={{this.path}}
            @backend={{this.backend}}
            @breadcrumbs={{this.breadcrumbs}}
        />`,
        { owner: this.engine }
      );

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

    assert.dom(FORM.dataInputLabel({ isJson: false })).hasText('Version data');
    assert.dom('.cm-editor').doesNotExist('CodeMirror editor is not rendered');

    await click(GENERAL.toggleInput('json'));
    assert.dom(FORM.dataInputLabel({ isJson: true })).hasText('Version data');
    await waitFor('.cm-editor');
    assert.dom('.cm-editor').exists('CodeMirror editor is rendered');
  });

  test('it should show old version alert', async function (assert) {
    this.metadata.current_version = 2;
    await this.renderComponent();
    assert
      .dom(FORM.versionAlert)
      .hasText(
        `Warning You are creating a new version based on data from Version 1. The current version for my-secret is Version 2.`
      );
  });

  test('it should render fail read error', async function (assert) {
    this.secret.failReadErrorCode = 403;
    await this.renderComponent();
    assert.dom(FORM.noReadAlert).exists('it renders no read alert');
  });

  test('it should render diff view', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.toggleInput('Show diff')).isNotDisabled('Diff toggle is not disabled');
    assert.dom(PAGE.edit.toggleDiffDescription).hasText('No changes to show. Update secret to view diff');
    assert.dom(PAGE.diff.visualDiff).doesNotExist('Does not show visual diff');

    await fillIn(FORM.keyInput(1), 'foo2');
    await fillIn(FORM.maskedValueInput(1), 'bar2');

    assert.dom(GENERAL.toggleInput('Show diff')).isNotDisabled('Diff toggle is not disabled');
    assert.dom(PAGE.edit.toggleDiffDescription).hasText('Showing the diff will reveal secret values');
    assert.dom(PAGE.diff.visualDiff).doesNotExist('Does not show visual diff');
    await click(GENERAL.toggleInput('Show diff'));
    assert.dom(PAGE.diff.visualDiff).exists('Shows visual diff');
    assert.dom(PAGE.diff.added).hasText(`foo2"bar2"`);

    await click(GENERAL.toggleInput('json'));
    await waitFor('.cm-editor');
    const editor = codemirror();

    setCodeEditorValue(editor, '{ "foo3": "bar3" }');
    await settled();

    assert.dom(PAGE.diff.visualDiff).exists('Visual diff updates');
    assert.dom(PAGE.diff.deleted).hasText(`foo"bar"`);
    assert.dom(PAGE.diff.added).hasText(`foo3"bar3"`);
  });
});
