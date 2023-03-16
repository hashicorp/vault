import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | sidebar-user-menu', function (hooks) {
  setupRenderingTest(hooks);

  test('it should render trigger and open menu', async function (assert) {
    await render(hbs`<Sidebar::UserMenu />`);
    assert
      .dom('[data-test-user-menu-trigger] [data-test-icon="user"]')
      .exists('Correct icon renders for menu trigger');
    await click('[data-test-user-menu-trigger]');
    assert.dom('[data-test-user-menu-content]').exists('User menu content renders');
  });

  test('it should render default menu items', async function (assert) {
    const auth = this.owner.lookup('service:auth');
    sinon.stub(auth, 'currentToken').value('root');

    await render(hbs`<Sidebar::UserMenu />`);
    await click('[data-test-user-menu-trigger]');

    assert.dom('.menu-label').hasText('Token', 'Auth data display name renders');
    assert.dom('li').exists({ count: 2 }, 'Correct number of menu items render');
    assert.dom('[data-clipboard-text="root"]').exists('Copy token action renders');
    assert.dom('#logout').hasText('Log out', 'Log out action renders');
  });

  test('it should render conditional menu items', async function (assert) {
    const router = this.owner.lookup('service:router');
    const transitionStub = sinon.stub(router, 'transitionTo');
    const auth = this.owner.lookup('service:auth');
    const renewStub = sinon.stub(auth, 'renew').resolves();
    const revokeStub = sinon.stub(auth, 'revokeCurrentToken').resolves();
    const date = new Date();
    sinon.stub(auth, 'tokenExpirationDate').value(date.setDate(date.getDate() + 1));
    sinon.stub(auth, 'authData').value({ displayName: 'token', renewable: true, entity_id: 'foo' });
    auth.set('allowExpiration', true);

    await render(hbs`<Sidebar::UserMenu />`);
    await click('[data-test-user-menu-trigger]');

    assert.dom('[data-test-user-menu-item="token alert"]').exists('Token expiration alert renders');
    assert.dom('[data-test-user-menu-item="mfa"]').hasText('Multi-factor authentication', 'MFA link renders');

    await click('[data-test-user-menu-item="revoke token"]');
    await click('[data-test-confirm-button]');
    assert.true(revokeStub.calledOnce, 'Auth revoke token method called on revoke confirm');
    assert.true(
      transitionStub.calledWith('vault.cluster.logout'),
      'Route transitions to log out on revoke success'
    );

    await click('[data-test-user-menu-item="renew token"]');
    assert.true(renewStub.calledOnce, 'Auth renew token method called');
  });
});
