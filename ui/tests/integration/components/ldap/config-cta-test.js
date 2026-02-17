/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | ldap | ConfigCta', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'ldap');
  setupMirage(hooks);

  test('it should render message and action', async function (assert) {
    await render(hbs`<ConfigCta />`, { owner: this.engine });
    assert.dom(GENERAL.emptyStateTitle).hasText('LDAP not configured', 'Title renders');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText('Get started by setting up the connection with your existing LDAP system.', 'Message renders');
    assert.dom(GENERAL.emptyStateActions).hasText('Configure LDAP', 'Action renders');
  });
});
