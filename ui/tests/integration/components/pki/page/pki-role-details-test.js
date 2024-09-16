/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { PKI_ROLE_DETAILS } from 'vault/tests/helpers/pki/pki-selectors';

module('Integration | Component | pki role details page', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/role', {
      name: 'Foobar',
      backend: 'pki',
      noStore: false,
      noStoreMetadata: true,
      keyUsage: [],
      extKeyUsage: ['bar', 'baz'],
      ttl: 600,
    });
  });

  test('it should render the page component', async function (assert) {
    await render(
      hbs`
      <Page::PkiRoleDetails @role={{this.model}} />
  `,
      { owner: this.engine }
    );
    assert.dom(PKI_ROLE_DETAILS.issuerLabel).hasText('Issuer', 'Label is');
    assert
      .dom(`${PKI_ROLE_DETAILS.keyUsageValue} [data-test-icon="minus"]`)
      .exists('Key usage shows dash when array is empty');
    assert
      .dom(PKI_ROLE_DETAILS.extKeyUsageValue)
      .hasText('bar,baz', 'Key usage shows comma-joined values when array has items');
    assert
      .dom(PKI_ROLE_DETAILS.noStoreValue)
      .containsText('Yes', 'noStore shows opposite of what the value is');
    assert
      .dom(PKI_ROLE_DETAILS.noStoreMetadataValue)
      .doesNotExist('does not render value for enterprise-only field');
    assert.dom(PKI_ROLE_DETAILS.customTtlValue).containsText('10 minutes', 'TTL shown as duration');
  });

  test('it should render the enterprise-only values in enterprise edition', async function (assert) {
    const version = this.owner.lookup('service:version');
    version.type = 'enterprise';
    await render(
      hbs`
      <Page::PkiRoleDetails @role={{this.model}} />
  `,
      { owner: this.engine }
    );
    assert
      .dom(PKI_ROLE_DETAILS.noStoreMetadataValue)
      .containsText('No', 'noStoreMetadata shows opposite of what the value is');
  });

  test('it should render the notAfter date if present', async function (assert) {
    assert.expect(1);
    this.model = this.store.createRecord('pki/role', {
      name: 'Foobar',
      backend: 'pki',
      noStore: false,
      keyUsage: [],
      extKeyUsage: ['bar', 'baz'],
      notAfter: '2030-05-04T12:00:00.000Z',
    });
    await render(
      hbs`
      <Page::PkiRoleDetails @role={{this.model}} />
  `,
      { owner: this.engine }
    );
    assert
      .dom(PKI_ROLE_DETAILS.customTtlValue)
      .containsText('May', 'Formats the notAfter date instead of TTL');
  });
});
