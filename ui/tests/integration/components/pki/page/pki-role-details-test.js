/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/page/pki-role-details';

module('Integration | Component | pki role details page', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/role', {
      name: 'Foobar',
      backend: 'pki',
      noStore: false,
      keyUsage: [],
      extKeyUsage: ['bar', 'baz'],
      ttl: 600,
    });
  });

  test('it should render the page component', async function (assert) {
    assert.expect(5);
    await render(
      hbs`
      <Page::PkiRoleDetails @role={{this.model}} />
  `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.issuerLabel).hasText('Issuer', 'Label is');
    assert.dom(SELECTORS.keyUsageValue).hasText('None', 'Key usage shows none when array is empty');
    assert
      .dom(SELECTORS.extKeyUsageValue)
      .hasText('bar, baz,', 'Key usage shows comma-joined values when array has items');
    assert.dom(SELECTORS.noStoreValue).containsText('Yes', 'noStore shows opposite of what the value is');
    assert.dom(SELECTORS.customTtlValue).containsText('10m', 'TTL shown as duration');
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
    assert.dom(SELECTORS.customTtlValue).containsText('May', 'Formats the notAfter date instead of TTL');
  });
});
