/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import { CreationMethod } from 'vault/utils/constants/snippet';
import sinon from 'sinon';

module('Integration | Component | pki | external-pki | ExternalPki::Page::Overview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    const router = this.owner.lookup('service:router');
    this.transitionToStub = sinon.stub(router, 'transitionTo');
    this.setError = (status) => {
      let message = '';
      switch (status) {
        case 404:
          message = '';
          break;
        case 403:
          message = '1 error occurred:\n\t* permission denied\n\n';
          break;
        case 500:
          message = 'Internal server error';
          break;
        default:
          // This is what is the route model returns when there is no error
          return { message: '' };
      }
      return { status, message };
    };
    this.setupModel = (respOverrides = {}) => {
      return {
        engine: new SecretsEngineResource({
          accessor: 'pki-external-ca_e158c567',
          type: 'pki-external-ca',
          path: 'my-pki-external-ca/',
        }),
        acmeAccountsResp: { keys: [], error: this.setError(404) },
        dnsProvidersResp: { keys: [], error: this.setError(404) },
        rolesResp: { keys: [], error: this.setError(404) },
        showConfigSnippets: undefined,
        ...respOverrides,
      };
    };

    this.renderComponent = () =>
      render(
        hbs`<ExternalPki::Page::Overview @model={{this.model}} />
        `,
        { owner: this.engine }
      );
  });

  test('it hides overview cards when user does not have permission', async function (assert) {
    this.model = this.setupModel({
      acmeAccountsResp: { keys: [], error: this.setError(403) },
      dnsProvidersResp: { keys: [], error: this.setError(403) },
      rolesResp: { keys: [], error: this.setError(403) },
    });
    await this.renderComponent();

    // Implementation select should not render
    assert.dom('h1').doesNotExist();
    assert.dom(GENERAL.radioCardByAttr()).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('terraform')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('cli')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('api')).doesNotExist();

    assert.dom(GENERAL.overviewCard.container('ACME accounts')).doesNotExist();
    assert.dom(GENERAL.overviewCard.container('DNS providers')).doesNotExist();
    assert.dom(GENERAL.overviewCard.container('Roles')).doesNotExist();

    // Only order and cert lookups render
    assert.dom(GENERAL.overviewCard.container('View certificate')).exists();
    assert.dom(GENERAL.overviewCard.container('Orders')).exists();
  });

  test('it renders implementation select when showConfigSnippets is true', async function (assert) {
    this.model = this.setupModel({ showConfigSnippets: true });
    await this.renderComponent();

    assert.dom('h1').hasText('Choose your implementation method');
    assert
      .dom(GENERAL.radioCardByAttr(CreationMethod.TERRAFORM))
      .exists()
      .isChecked('Terraform is initially selected');
    assert.dom(GENERAL.radioCardByAttr(CreationMethod.APICLI)).exists().isNotChecked();

    assert
      .dom(GENERAL.fieldByAttr('terraform'))
      .exists({ count: 2 }, 'Terraform snippet renders for each config');
    assert.dom(GENERAL.fieldByAttr('cli')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('api')).doesNotExist();
  });

  module('is configured', function (hooks) {
    hooks.beforeEach(function () {
      this.model = this.setupModel({
        acmeAccountsResp: { keys: ['account-1', 'account-2'], error: this.setError() },
        dnsProvidersResp: { keys: [], error: this.setError(404) },
        rolesResp: { keys: ['role-1', 'role-2', 'role-3'], error: this.setError() },
        showConfigSnippets: false,
      });
    });

    test('it renders overview cards with zero and non-zero counts', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.content('ACME accounts')).hasText('2');
      assert.dom(GENERAL.overviewCard.content('DNS providers')).hasText('0');
      assert.dom(GENERAL.overviewCard.content('Roles')).hasText('3');
    });

    test('it hides only ACME accounts card', async function (assert) {
      this.model = this.setupModel({ acmeAccountsResp: { keys: [], error: this.setError(403) } });
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.container('ACME accounts')).doesNotExist();
      assert.dom(GENERAL.overviewCard.container('DNS providers')).exists();
      assert.dom(GENERAL.overviewCard.container('Roles')).exists();
    });

    test('it hides only DNS providers card', async function (assert) {
      this.model = this.setupModel({ dnsProvidersResp: { keys: [], error: this.setError(403) } });
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.container('ACME accounts')).exists();
      assert.dom(GENERAL.overviewCard.container('DNS providers')).doesNotExist();
      assert.dom(GENERAL.overviewCard.container('Roles')).exists();
    });

    test('it hides only Roles card', async function (assert) {
      this.model = this.setupModel({ rolesResp: { keys: [], error: this.setError(403) } });
      await this.renderComponent();
      assert.dom(GENERAL.overviewCard.container('ACME accounts')).exists();
      assert.dom(GENERAL.overviewCard.container('DNS providers')).exists();
      assert.dom(GENERAL.overviewCard.container('Roles')).doesNotExist();
    });

    test('it displays error message for non-403 error', async function (assert) {
      this.model = this.setupModel({
        acmeAccountsResp: { keys: ['myaccount'], error: this.setError(500) },
        dnsProvidersResp: { keys: ['mydns'], error: this.setError(500) },
        rolesResp: { keys: ['myrole'], error: this.setError(500) },
      });
      await this.renderComponent();
      assert
        .dom(`${GENERAL.overviewCard.container('ACME accounts')} ${GENERAL.messageError}`)
        .exists()
        .hasText('Error Internal server error');
      assert.dom(GENERAL.overviewCard.content('ACME accounts')).doesNotExist();
      assert
        .dom(`${GENERAL.overviewCard.container('DNS providers')} ${GENERAL.messageError}`)
        .exists()
        .hasText('Error Internal server error');
      assert.dom(GENERAL.overviewCard.content('DNS providers')).doesNotExist();
      assert
        .dom(`${GENERAL.overviewCard.container('Roles')} ${GENERAL.messageError}`)
        .exists()
        .hasText('Error Internal server error');
      assert.dom(GENERAL.overviewCard.content('Roles')).doesNotExist();
    });

    test('it transitions to look up certificate', async function (assert) {
      await this.renderComponent();
      await fillIn(`${GENERAL.overviewCard.container('View certificate')} input`, '03:e7:1f:');
      await click(GENERAL.button('Lookup certificate'));
      const [route, engineId, serialNumber] = this.transitionToStub.lastCall.args;
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.pki.external.certificates.certificate',
        'it calls transition with expected route'
      );
      assert.strictEqual(engineId, 'my-pki-external-ca', 'it calls transition with expected engine ID');
      assert.strictEqual(serialNumber, '03:e7:1f:', 'it calls transition with inputted serial number');
    });

    test('it transitions to look up order', async function (assert) {
      await this.renderComponent();
      await fillIn(`${GENERAL.overviewCard.container('Orders')} input`, '01936d8e-7c3');
      await click(GENERAL.button('Lookup order'));
      const [route, engineId, orderId] = this.transitionToStub.lastCall.args;
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.pki.external.orders.order',
        'it calls transition with expected route'
      );
      assert.strictEqual(engineId, 'my-pki-external-ca', 'it calls transition with expected engine ID');
      assert.strictEqual(orderId, '01936d8e-7c3', 'it calls transition with inputted order ID');
    });
  });
});
