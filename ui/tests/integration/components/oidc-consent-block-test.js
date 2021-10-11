import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';

const redirectBase = 'https://hashicorp.com';

module('Integration | Component | oidc-consent-block', function(hooks) {
  setupRenderingTest(hooks);
  this.set('redirect', redirectBase);

  test('it renders', async function(assert) {
    await render(hbs`
      <OidcConsentBlock @redirect={{redirect}} @code="1234" />
    `);

    assert.dom('[data-test-consent-title]').hasText('Consent', 'Title is correct on initial render');
    assert
      .dom('[data-test-consent-form]')
      .includesText(
        'In order to complete the login process, you must consent to Vault sharing your profile, email, address, and phone with the client.',
        'shows the correct copy for consent form'
      );
    assert.dom('[data-test-edit-form-submit]').hasText('Yes', 'form button has correct submit text');
    assert.dom('[data-test-cancel-button]').hasText('No', 'form button has correct cancel text');
  });

  test('it calls the success callback when user clicks "Yes"', async function(assert) {
    const spy = sinon.spy();
    this.set('successSpy', spy);
    this.set('redirect', redirectBase);

    await render(hbs`
      <OidcConsentBlock @redirect={{redirectBase}} @code="1234" @onSuccess={{successSpy}} />
    `);

    assert.dom('[data-test-consent-title]').hasText('Consent', 'Title is correct on initial render');
    assert.dom('[data-test-consent-form]').exists('Consent form exists');
    assert
      .dom('[data-test-consent-form]')
      .includesText(
        'In order to complete the login process, you must consent to Vault sharing your profile, email, address, and phone with the client.',
        'shows the correct copy for consent form'
      );
    await click('[data-test-edit-form-submit]');
    assert.ok(
      spy.calledWith(redirectBase, { code: 1234 }),
      'calls the radio change function when option clicked'
    );
  });

  test('it shows the termination message when user clicks "No"', async function(assert) {
    const spy = sinon.spy();
    this.set('successSpy', spy);
    this.set('redirect', redirectBase);

    await render(hbs`
      <OidcConsentBlock @redirect={{redirectBase}} @code="1234" @onSuccess={{successSpy}} />
    `);

    assert.dom('[data-test-consent-title]').hasText('Consent', 'Title is correct on initial render');
    assert.dom('[data-test-consent-form]').exists('Consent form exists');
    assert
      .dom('[data-test-consent-form]')
      .includesText(
        'In order to complete the login process, you must consent to Vault sharing your profile, email, address, and phone with the client.',
        'shows the correct copy for consent form'
      );
    await click('[data-test-cancel-button]');
    assert.dom('[data-test-consent-title]').hasText('Consent Not Given', 'Title changes to not given');
    assert.dom('[data-test-consent-form]').doesNotExist('Consent form is hidden');

    assert.ok(spy.notCalled, 'Does not call the success method');
  });
});
