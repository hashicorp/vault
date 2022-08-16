import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
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
    await render(hbs`<MountAccessorSelect @onChange={{onChange}}/>`);
    assert.dom('[data-test-mount-accessor-select]').exists();
  });

  test('it filters token', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{onChange}} @filterToken={{true}}/>`);
    await click('[data-test-mount-accessor-select]');
    let options = document.querySelector('[data-test-mount-accessor-select]').options;
    assert.equal(options.length, 1, 'only the auth option, no token');
  });

  test('it shows token', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{onChange}}/>`);
    await click('[data-test-mount-accessor-select]');
    let options = document.querySelector('[data-test-mount-accessor-select]').options;
    assert.equal(options.length, 2, 'both auth and token show');
  });

  test('it sends value to parent onChange', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{onChange}}/>`);
    await fillIn('[data-test-mount-accessor-select]', 'auth_userpass_1234');
    assert.ok(
      this.onChange.calledWith('auth_userpass_1234'),
      'Passes the auth method selected to the parent'
    );
  });

  test('it selects the first option if no default', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{onChange}} />`);
    let defaultSelection = document.querySelector('[data-test-mount-accessor-select]').options[0].innerHTML;
    // remove all non letters
    assert.equal(defaultSelection.replace(/\W/g, ''), 'userpassuserpass');
  });

  test('it shows Select one if yes default', async function (assert) {
    await render(hbs`<MountAccessorSelect @onChange={{onChange}} @noDefault={{true}} />`);
    let defaultSelection = document.querySelector('[data-test-mount-accessor-select]').options[0].innerHTML;
    // remove all non letters
    assert.equal(defaultSelection.replace(/\W/g, ''), 'Selectone');
  });
});
