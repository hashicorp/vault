/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled, select } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { hbs } from 'ember-cli-htmlbars';
import { typeInSearch, clickTrigger } from 'ember-power-select/test-support/helpers';
import searchSelect from '../../../pages/components/search-select';
import { setupMirage } from 'ember-cli-mirage/test-support';

const SELECTORS = {
  form: '[data-test-keymgmt-distribution-form]',
  keySection: '[data-test-keymgmt-dist-key]',
  keyTypeSection: '[data-test-keymgmt-dist-keytype]',
  providerInput: '[data-test-keymgmt-dist-provider]',
  operationsSection: '[data-test-keymgmt-dist-operations]',
  protectionsSection: '[data-test-keymgmt-dist-protections]',
  errorKey: '[data-test-keymgmt-error="key"]',
  errorNewKey: '[data-test-keymgmt-error="new-key"]',
  errorProvider: '[data-test-keymgmt-error="provider"]',
  inlineError: '[data-test-keymgmt-error]',
};

const ssComponent = create(searchSelect);

module('Integration | Component | keymgmt/distribute', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.set('backend', 'keymgmt');
    this.set('providers', ['provider-aws', 'provider-gcp', 'provider-azure']);
    this.server.get('/keymgmt/key', () => {
      return {
        data: {
          keys: ['example-1', 'example-2', 'example-3'],
        },
      };
    });
    this.server.get('/keymgmt/key/:name', (_, req) => {
      const name = req.params.name;
      return {
        data: {
          name,
          type: 'aes256-gcm96', // incompatible with azurekeyvault only
        },
      };
    });
    this.server.get('/keymgmt/kms/:name', (_, req) => {
      const name = req.params.name;
      let provider;
      switch (name) {
        case 'provider-aws':
          provider = 'awskms';
          break;
        case 'provider-azure':
          provider = 'azurekeyvault';
          break;
        default:
          provider = 'gcpckms';
          break;
      }
      return {
        data: {
          name,
          provider,
        },
      };
    });
    this.server.get('/keymgmt/kms', () => {
      return {
        data: {
          keys: ['provider-aws', 'provider-azure', 'provider-gcp'],
        },
      };
    });
  });

  test('it does not allow operation selection until valid key/provider combo selected', async function (assert) {
    assert.expect(6);
    await render(
      hbs`<Keymgmt::Distribute @backend="keymgmt" @key="example-1" @providers={{this.providers}} @onClose={{fn (mut this.onClose)}} />`
    );
    assert.dom(SELECTORS.operationsSection).hasAttribute('disabled');
    // Select
    await clickTrigger();
    assert.strictEqual(ssComponent.options.length, 3, 'shows all provider options');
    await typeInSearch('aws');
    await ssComponent.selectOption();
    await settled();
    assert.dom(SELECTORS.operationsSection).doesNotHaveAttribute('disabled');
    // Remove selection
    await ssComponent.deleteButtons.objectAt(0).click();
    // Select Azure
    await clickTrigger();
    await typeInSearch('azure');
    await ssComponent.selectOption();
    // await select(SELECTORS.providerInput, 'provider-azure');
    assert.dom(SELECTORS.operationsSection).hasAttribute('disabled');
    assert.dom(SELECTORS.inlineError).exists({ count: 1 }, 'only shows single error');
    assert.dom(SELECTORS.errorProvider).exists('Shows key/provider match error on provider');
  });
  test('it shows key type select field if new key created', async function (assert) {
    assert.expect(2);
    await render(
      hbs`<Keymgmt::Distribute @backend="keymgmt" @providers={{this.providers}} @onClose={{fn (mut this.onClose)}} />`
    );
    assert.dom(SELECTORS.keyTypeSection).doesNotExist('Key Type section is not rendered by default');
    // Add new item on search-select
    await clickTrigger();
    await typeInSearch('new-key');
    await ssComponent.selectOption();
    assert.dom(SELECTORS.keyTypeSection).exists('Key Type selector is shown');
  });
  test('it hides the provider field if passed from the parent', async function (assert) {
    assert.expect(5);
    await render(
      hbs`<Keymgmt::Distribute @backend="keymgmt" @provider="provider-azure" @onClose={{fn (mut this.onClose)}} />`
    );
    assert.dom(SELECTORS.providerInput).doesNotExist('Provider input is hidden');
    // Select existing key
    await clickTrigger();
    await ssComponent.selectOption();
    await settled();
    assert.dom(SELECTORS.inlineError).exists({ count: 1 }, 'only shows single error');
    assert.dom(SELECTORS.errorKey).exists('Shows error on key selector when key/provider mismatch');
    // Remove selection
    await ssComponent.deleteButtons.objectAt(0).click();
    // Select new key
    await clickTrigger();
    await typeInSearch('new-key');
    await ssComponent.selectOption();
    await select(SELECTORS.keyTypeSection, 'ecdsa-p256');
    assert.dom(SELECTORS.inlineError).exists({ count: 1 }, 'only shows single error');
    assert.dom(SELECTORS.errorNewKey).exists('Shows error on key type');
  });
  test('it hides the key field if passed from the parent', async function (assert) {
    assert.expect(4);
    await render(
      hbs`<Keymgmt::Distribute @backend="keymgmt" @providers={{this.providers}} @key="example-1" @onClose={{fn (mut this.onClose)}} />`
    );
    assert.dom(SELECTORS.providerInput).exists('Provider input shown');
    assert.dom(SELECTORS.keySection).doesNotExist('Key input not shown');
    await clickTrigger();
    await typeInSearch('azure');
    await ssComponent.selectOption();
    assert.dom(SELECTORS.inlineError).exists({ count: 1 }, 'only shows single error');
    assert.dom(SELECTORS.errorProvider).exists('Shows error due to key/provider mismatch');
  });
});
