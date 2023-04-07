/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/page/pki-tidy-form';

module('Integration | Component | pki | Page::PkiTidyForm', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';

    this.store.createRecord('pki/tidy', { backend: 'pki-test' });

    this.tidy = this.store.peekAll('pki/tidy');
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: 'pki-test', route: 'overview' },
      { label: 'configuration', route: 'configuration.index' },
      { label: 'tidy' },
    ];
  });

  test('it should render tidy fields', async function (assert) {
    await render(hbs`<Page::PkiTidyForm @tidy={{this.tidy}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    assert.dom(SELECTORS.tidyCertStoreLabel).hasText('Tidy the certificate store');
    assert.dom(SELECTORS.tidyRevocationList).hasText('Tidy the revocation list (CRL)');
    assert.dom(SELECTORS.safetyBufferTTL).exists();
  });

  test('it should change the attributes on the model', async function (assert) {
    await render(hbs`<Page::PkiTidyForm @tidy={{this.tidy}} @breadcrumbs={{this.breadcrumbs}} />`, {
      owner: this.engine,
    });
    await click(SELECTORS.tidyCertStoreCheckbox);
    await click(SELECTORS.tidyRevocationCheckbox);
    await fillIn(SELECTORS.safetyBufferInput, '72h');
    assert.true(this.tidy.tidyCertStore);
    assert.true(this.tidy.tidyRevocationQueue);
    assert.dom(SELECTORS.safetyBufferInput).hasValue('72h');
  });
});
