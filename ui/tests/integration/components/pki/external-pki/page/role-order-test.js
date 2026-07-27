/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | pki | external-pki | ExternalPki::Page::RolesRoleOrder', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');

    this.model = {
      engine: { id: 'pki-external-ca' },
      order_id: 'order-abc-123',
      order: { details: { order_status: 'pending' }, error: undefined },
      certificate: { details: undefined, error: undefined },
      responseTimestamp: new Date('2026-07-14T21:00:00Z'),
    };

    this.renderComponent = () =>
      render(hbs`<ExternalPki::Page::RolesRoleOrder @model={{this.model}} />`, { owner: this.engine });
  });

  test('it renders the last refreshed timestamp', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.textBody('Last refreshed')).hasTextContaining('Last refreshed: July 14, 2026');
  });

  test('it calls router.refresh with the role order route when Refresh button is clicked', async function (assert) {
    const refreshStub = sinon.stub(this.router, 'refresh');
    await this.renderComponent();
    await click(GENERAL.button('Refresh'));
    assert.true(refreshStub.calledOnce, 'refresh was called once');
    assert.true(
      refreshStub.calledWith('vault.cluster.secrets.backend.pki.external.roles.role.order'),
      'refresh was called with the role order route'
    );
  });
});
