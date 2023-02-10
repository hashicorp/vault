import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, typeIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';

module('Integration | Component | mount-accessor-select', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.server.get('/sys/auth', () => ({
      data: {
        'userpass/': { type: 'userpass', accessor: 'auth_userpass_1234' },
        'token/': { type: 'token', accessor: 'auth_token' },
      },
    }));
    this.set('onChange', sinon.spy());
  });

  test('it renders', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{this.onChange}}/>`);
    assert.dom('[data-test-mount-accessor-select]').exists();
  });

  test('it filters token', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{this.onChange}} @filterToken={{true}}/>`);
    await click('[data-test-mount-accessor-select]');
    const options = document.querySelector('[data-test-mount-accessor-select]').options;
    assert.strictEqual(options.length, 1, 'only the auth option, no token');
  });

  test('it shows token', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{this.onChange}}/>`);
    await click('[data-test-mount-accessor-select]');
    const options = document.querySelector('[data-test-mount-accessor-select]').options;
    assert.strictEqual(options.length, 2, 'both auth and token show');
  });

  test('it sends value to parent onChange', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{this.onChange}}/>`);
    await typeIn('[data-test-mount-accessor-select]', 'auth_userpass_1234');
    assert.ok(
      this.onChange.calledWith('auth_userpass_1234'),
      'Passes the auth method selected to the parent'
    );
  });

  test('it selects the first option if no default', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{this.onChange}} />`);
    const defaultSelection = document.querySelector('[data-test-mount-accessor-select]').options[0].innerHTML;
    // remove all non letters
    assert.strictEqual(defaultSelection.replace(/\W/g, ''), 'userpassuserpass');
  });

  test('it shows Select one if yes default', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{this.onChange}} @noDefault={{true}} />`);
    const defaultSelection = document.querySelector('[data-test-mount-accessor-select]').options[0].innerHTML;
    // remove all non letters
    assert.strictEqual(defaultSelection.replace(/\W/g, ''), 'Selectone');
  });
});
