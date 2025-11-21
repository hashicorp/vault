/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | sidebar-user-menu', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.auth = this.owner.lookup('service:auth');
    setRunOptions({
      // TODO: fix this component
      rules: {
        'nested-interactive': { enabled: false },
        // TODO: fix ConfirmAction rendered in toolbar not a list item
        list: { enabled: false },
      },
    });
  });

  test('it should render trigger and open menu', async function (assert) {
    await render(hbs`<Sidebar::UserMenu />`);

    assert
      .dom(`${GENERAL.button('user-menu-trigger')} [data-test-icon="user"]`)
      .exists('Correct icon renders for menu trigger');
    await click(GENERAL.button('user-menu-trigger'));
    assert.dom('[data-test-user-menu-item="title"]').exists('User menu content renders');
  });

  test('it should render default menu items', async function (assert) {
    sinon.stub(this.auth, 'currentToken').value('root');
    sinon.stub(this.auth, 'authData').value({ displayName: 'token' });

    await render(hbs`<Sidebar::UserMenu />`);
    await click(GENERAL.button('user-menu-trigger'));

    assert.dom('[data-test-user-menu-item="title"]').hasText('Token', 'Auth data display name renders');
    assert.dom('li').exists({ count: 3 }, 'Correct number of menu items render');
    assert.dom(GENERAL.copyButton).exists('Copy token action renders');
    assert.dom('[data-test-user-menu-item="logout"]').hasText('Log out', 'Log out action renders');
  });

  test('it should render conditional menu items', async function (assert) {
    const date = new Date();
    sinon.stub(this.auth, 'tokenExpirationDate').value(date.setDate(date.getDate() + 1));
    sinon.stub(this.auth, 'authData').value({ displayName: 'token', renewable: true, entityId: 'foo' });
    this.auth.set('allowExpiration', true);

    await render(hbs`<Sidebar::UserMenu />`);
    await click(GENERAL.button('user-menu-trigger'));

    assert.dom('[data-test-user-menu-item="token alert"]').exists('Token expiration alert renders');
    assert.dom('[data-test-user-menu-item="mfa"]').hasText('Multi-factor authentication', 'MFA link renders');

    assert.dom('[data-test-user-menu-item="renew token"]').exists('Renew token action renders');
    assert.dom('[data-test-user-menu-item="revoke token"]').exists('Revoke token action renders');
  });

  test('it should renew token', async function (assert) {
    const date = new Date();
    const renewStub = sinon.stub(this.auth, 'renew').resolves();
    sinon.stub(this.auth, 'tokenExpirationDate').value(date.setDate(date.getDate() + 1));
    sinon.stub(this.auth, 'authData').value({ displayName: 'token', renewable: true, entityId: 'foo' });

    await render(hbs`<Sidebar::UserMenu />`);
    await click(GENERAL.button('user-menu-trigger'));

    await click('[data-test-user-menu-item="renew token"]');
    assert.true(renewStub.calledOnce, 'Auth renew token method called');
  });

  test('it should revoke token', async function (assert) {
    const date = new Date();
    const router = this.owner.lookup('service:router');
    const transitionStub = sinon.stub(router, 'transitionTo');
    const revokeStub = sinon.stub(this.auth, 'revokeCurrentToken').resolves();
    sinon.stub(this.auth, 'tokenExpirationDate').value(date.setDate(date.getDate() + 1));
    sinon.stub(this.auth, 'authData').value({ displayName: 'token', entityId: 'foo' });

    await render(hbs`<Sidebar::UserMenu />`);
    await click(GENERAL.button('user-menu-trigger'));

    await click('[data-test-user-menu-item="revoke token"]');
    await click(GENERAL.confirmButton);
    assert.true(revokeStub.calledOnce, 'Auth revoke token method called on revoke confirm');
    assert.true(
      transitionStub.calledWith('vault.cluster.logout'),
      'Route transitions to log out on revoke success'
    );
  });
});
