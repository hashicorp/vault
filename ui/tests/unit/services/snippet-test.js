/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { CreationMethod } from 'vault/utils/constants/snippet';

const TF_SNIPPET = 'resource "vault_mount" "example" {}';
const API_SNIPPET = 'curl --header "X-Vault-Token: $VAULT_TOKEN" $VAULT_ADDR/v1/test';
const CLI_SNIPPET = 'vault write test/path key=value';

const CUSTOM_TABS = [
  { key: 'api', label: 'API', snippet: API_SNIPPET },
  { key: 'cli', label: 'CLI', snippet: CLI_SNIPPET },
];

module('Unit | Service | snippet', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.service = this.owner.lookup('service:snippet');
  });

  module('default state', function () {
    test('selectedTabIdx defaults to 0', function (assert) {
      assert.strictEqual(this.service.selectedTabIdx, 0);
    });

    test('creationMethodChoice defaults to TERRAFORM', function (assert) {
      assert.strictEqual(this.service.creationMethodChoice, CreationMethod.TERRAFORM);
    });

    test('codeSnippet defaults to null', function (assert) {
      assert.strictEqual(this.service.codeSnippet, null);
    });
  });

  module('#reset', function () {
    test('resets service state to defaults', function (assert) {
      this.service.selectedTabIdx = 1;
      this.service.creationMethodChoice = CreationMethod.APICLI;
      this.service.codeSnippet = TF_SNIPPET;

      this.service.reset();

      assert.strictEqual(this.service.selectedTabIdx, 0);
      assert.strictEqual(this.service.creationMethodChoice, CreationMethod.TERRAFORM);
      assert.strictEqual(this.service.codeSnippet, null);
    });

    test('accepts an initial creation method', function (assert) {
      this.service.reset(CreationMethod.APICLI);

      assert.strictEqual(this.service.creationMethodChoice, CreationMethod.APICLI);
      assert.strictEqual(this.service.selectedTabIdx, 0);
      assert.strictEqual(this.service.codeSnippet, null);
    });
  });

  module('#persistSnippet', function () {
    test('sets codeSnippet to tfSnippet when method is TERRAFORM', function (assert) {
      this.service.creationMethodChoice = CreationMethod.TERRAFORM;

      this.service.persistSnippet(TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.codeSnippet, TF_SNIPPET);
    });

    test('sets codeSnippet to the selected tab snippet when method is APICLI', function (assert) {
      this.service.creationMethodChoice = CreationMethod.APICLI;
      this.service.selectedTabIdx = 0;

      this.service.persistSnippet(TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.codeSnippet, API_SNIPPET);
    });

    test('sets codeSnippet to the CLI tab when selectedTabIdx is 1', function (assert) {
      this.service.creationMethodChoice = CreationMethod.APICLI;
      this.service.selectedTabIdx = 1;

      this.service.persistSnippet(TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.codeSnippet, CLI_SNIPPET);
    });

    test('sets codeSnippet to null when method is UI', function (assert) {
      this.service.creationMethodChoice = CreationMethod.UI;

      this.service.persistSnippet(TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.codeSnippet, null);
    });

    test('sets codeSnippet to null when APICLI tab index is out of range', function (assert) {
      this.service.creationMethodChoice = CreationMethod.APICLI;
      this.service.selectedTabIdx = 99;

      this.service.persistSnippet(TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.codeSnippet, null);
    });
  });

  module('#setCreationMethod', function () {
    test('updates creationMethodChoice', function (assert) {
      this.service.setCreationMethod(CreationMethod.APICLI, TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.creationMethodChoice, CreationMethod.APICLI);
    });

    test('persists the correct snippet after changing to TERRAFORM', function (assert) {
      this.service.setCreationMethod(CreationMethod.TERRAFORM, TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.codeSnippet, TF_SNIPPET);
    });

    test('persists the correct snippet after changing to APICLI', function (assert) {
      this.service.selectedTabIdx = 0;
      this.service.setCreationMethod(CreationMethod.APICLI, TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.codeSnippet, API_SNIPPET);
    });

    test('clears codeSnippet when changing to UI', function (assert) {
      this.service.codeSnippet = TF_SNIPPET;
      this.service.setCreationMethod(CreationMethod.UI, TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.codeSnippet, null);
    });
  });

  module('#setSelectedTab', function () {
    test('updates selectedTabIdx', function (assert) {
      this.service.creationMethodChoice = CreationMethod.APICLI;

      this.service.setSelectedTab(1, TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.selectedTabIdx, 1);
    });

    test('persists the snippet for the newly selected tab', function (assert) {
      this.service.creationMethodChoice = CreationMethod.APICLI;

      this.service.setSelectedTab(1, TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.codeSnippet, CLI_SNIPPET);
    });

    test('persists tfSnippet if method is still TERRAFORM when tab changes', function (assert) {
      this.service.creationMethodChoice = CreationMethod.TERRAFORM;

      this.service.setSelectedTab(1, TF_SNIPPET, CUSTOM_TABS);

      assert.strictEqual(this.service.codeSnippet, TF_SNIPPET);
    });
  });
});
