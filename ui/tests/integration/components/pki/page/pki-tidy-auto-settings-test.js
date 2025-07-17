/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';

module('Integration | Component | page/pki-tidy-auto-settings', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    const backend = 'pki-auto-tidy';
    this.backend = backend;

    this.context = { owner: this.engine };
    this.store = this.owner.lookup('service:store');

    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: backend, route: 'overview', model: backend },
      { label: 'Tidy', route: 'tidy.index', model: backend },
      { label: 'Auto' },
    ];
  });

  test('it renders', async function (assert) {
    const model = this.store.createRecord('pki/tidy', {
      backend: this.backend,
      tidyType: 'auto',
      enabled: false,
      intervalDuration: '2d',
      tidyCertStore: false,
      tidyExpiredIssuers: true,
    });
    this.set('model', model);

    await render(
      hbs`<Page::PkiTidyAutoSettings @breadcrumbs={{this.breadcrumbs}} @model={{this.model}} />`,
      this.context
    );

    assert.dom('[data-test-breadcrumbs] li').exists({ count: 4 }, 'an item exists for each breadcrumb');
    assert.dom('[data-test-header-title]').hasText('Automatic Tidy Configuration', 'title is correct');
    assert
      .dom('[data-test-pki-edit-tidy-auto-link]')
      .hasText('Edit auto-tidy', 'toolbar edit link has correct text');

    assert.dom('[data-test-row="enabled"] [data-test-label-div]').hasText('Automatic tidy enabled');
    assert.dom('[data-test-value-div="Automatic tidy enabled"]').hasText('No');
    assert.dom('[data-test-value-div="Interval duration"]').hasText('2 days');
    // Universal operations
    assert.dom('[data-test-group-title="Universal operations"]').hasText('Universal operations');
    assert
      .dom('[data-test-value-div="Tidy the certificate store"]')
      .exists('Renders universal field when value exists');
    assert.dom('[data-test-value-div="Tidy the certificate store"]').hasText('No');
    assert
      .dom('[data-test-value-div="Tidy revoked certificates"]')
      .doesNotExist('Does not render universal field when value null');
    // Issuer operations
    assert.dom('[data-test-group-title="Issuer operations"]').hasText('Issuer operations');
    assert
      .dom('[data-test-value-div="Tidy expired issuers"]')
      .exists('Renders issuer op field when value exists');
    assert.dom('[data-test-value-div="Tidy expired issuers"]').hasText('Yes');
    assert
      .dom('[data-test-value-div="Tidy legacy CA bundle"]')
      .doesNotExist('Does not render issuer op field when value null');
  });
});
