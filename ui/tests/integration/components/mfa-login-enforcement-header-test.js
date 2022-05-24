import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | mfa-login-enforcement-header', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  test('it renders heading', async function (assert) {
    await render(hbs`<MfaLoginEnforcementHeader @heading="New enforcement" />`);

    assert.dom('[data-test-mleh-title]').includesText('New enforcement');
    assert.dom('[data-test-mleh-title] svg').hasClass('flight-icon-lock', 'Lock icon renders');
    assert
      .dom('[data-test-mleh-description]')
      .includesText('An enforcement will define which auth types', 'Description renders');
    assert.dom('[data-test-mleh-radio]').doesNotExist('Radio cards are hidden when not inline display mode');
    assert
      .dom('[data-test-component="search-select"]')
      .doesNotExist('Search select is hidden when not inline display mode');
  });

  test('it renders inline', async function (assert) {
    assert.expect(7);

    this.server.get('/identity/mfa/login-enforcement', () => {
      assert.ok(true, 'Request made to fetch enforcements');
      return {
        data: {
          key_info: {
            foo: { name: 'foo' },
          },
          keys: ['foo'],
        },
      };
    });

    await render(hbs`
      <MfaLoginEnforcementHeader
        @isInline={{true}}
        @radioCardGroupValue={{this.value}}
        @onRadioCardSelect={{fn (mut this.value)}}
        @onEnforcementSelect={{fn (mut this.enforcement)}}
      />
    `);

    assert.dom('[data-test-mleh-title]').includesText('Enforcement');
    assert
      .dom('[data-test-mleh-description]')
      .includesText('An enforcement includes the authentication types', 'Description renders');

    for (const option of ['new', 'existing', 'skip']) {
      await click(`[data-test-mleh-radio="${option}"] input`);
      assert.equal(this.value, option, 'Value is updated on radio select');
      if (option === 'existing') {
        await click('.ember-basic-dropdown-trigger');
        await click('.ember-power-select-option');
      }
    }

    assert.equal(this.enforcement.name, 'foo', 'Existing enforcement is selected');
  });
});
