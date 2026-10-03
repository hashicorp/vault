/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | page/pki-tidy-auto-settings', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.backend = 'pki-auto-tidy';
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'overview', model: this.backend },
      { label: 'Tidy', route: 'tidy.index', model: this.backend },
      { label: 'Auto' },
    ];

    this.model = {
      enabled: false,
      interval_duration: '2d',
      tidy_cert_store: false,
      tidy_expired_issuers: true,
    };

    this.renderComponent = () =>
      render(
        hbs`<Page::PkiTidyAutoSettings @breadcrumbs={{this.breadcrumbs}} @model={{this.model}} @backend={{this.backend}} />`,
        {
          owner: this.engine,
        }
      );
  });

  test('it renders', async function (assert) {
    await this.renderComponent();

    assert.dom('[data-test-breadcrumbs] li').exists({ count: 4 }, 'an item exists for each breadcrumb');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Automatic Tidy Configuration', 'title is correct');
    assert
      .dom('[data-test-pki-edit-tidy-auto-link]')
      .hasText('Edit auto-tidy', 'toolbar edit link has correct text');

    assert.dom('[data-test-row="enabled"] [data-test-row-label]').hasText('Automatic tidy enabled');
    assert.dom(GENERAL.infoRowValue('Automatic tidy enabled')).hasText('No');
    assert.dom(GENERAL.infoRowValue('Interval duration')).hasText('2 days');
    // Universal operations
    assert.dom('[data-test-group-title="Universal operations"]').hasText('Universal operations');
    assert
      .dom(GENERAL.infoRowValue('Tidy the certificate store'))
      .exists('Renders universal field when value exists');
    assert.dom(GENERAL.infoRowValue('Tidy the certificate store')).hasText('No');
    assert
      .dom(GENERAL.infoRowValue('Tidy revoked certificates'))
      .doesNotExist('Does not render universal field when value null');
    // Issuer operations
    assert.dom('[data-test-group-title="Issuer operations"]').hasText('Issuer operations');
    assert
      .dom(GENERAL.infoRowValue('Tidy expired issuers'))
      .exists('Renders issuer op field when value exists');
    assert.dom(GENERAL.infoRowValue('Tidy expired issuers')).hasText('Yes');
    assert
      .dom(GENERAL.infoRowValue('Tidy legacy CA bundle'))
      .doesNotExist('Does not render issuer op field when value null');
  });
});
