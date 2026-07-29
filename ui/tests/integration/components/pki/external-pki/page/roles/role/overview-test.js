/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render, select } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

const CERT_RESPONSE = {
  serial_number: 'ab:cd:ef:12:34:56',
  certificate: '-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----',
};

module(
  'Integration | Component | pki | external-pki | ExternalPki::Page::Roles::Role::Overview',
  function (hooks) {
    setupRenderingTest(hooks);
    setupEngine(hooks, 'pki');

    hooks.beforeEach(function () {
      this.model = {
        engine: { id: 'pki-external-ca' },
        role: {
          name: 'my-role',
          acme_account_name: 'production-account',
          dns_provider_name: 'aws-route53-prod',
          allowed_domains: ['example.com', '*.example.com'],
          allow_subdomains: true,
        },
      };

      const api = this.owner.lookup('service:api');
      this.fetchStub = sinon.stub(api.secrets, 'pkiExternalCaReadRoleCached').resolves(CERT_RESPONSE);
      // The component passes an addQueryParams callback as the third arg to pkiExternalCaReadRoleCached()
      this.addQueryParamsStub = sinon.stub(api, 'addQueryParams');

      this.renderComponent = () =>
        render(hbs`<ExternalPki::Page::Roles::Role::Overview @model={{this.model}} />`, {
          owner: this.engine,
        });
    });

    test('it renders role details', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.infoRowValue('ACME account name')).hasText('production-account');
      assert.dom(GENERAL.infoRowValue('DNS provider name')).hasText('aws-route53-prod');
      assert.dom(GENERAL.infoRowValue('Allowed domains')).hasText('example.com,*.example.com');
      assert.dom(GENERAL.infoRowValue('Allow subdomains')).hasText('Yes');

      assert
        .dom(GENERAL.overviewCard.container('Retrieve cached certificate'))
        .exists('fetch form is shown before a cert is fetched');
    });

    test('it does not render the cert details view before a successful fetch', async function (assert) {
      await this.renderComponent();

      // OrderCertDetails renders an info table with the certificate data — it should not appear yet
      assert.dom(GENERAL.overviewCard.container('Retrieve cached certificate')).exists('form is present');
      assert.dom(GENERAL.infoRowValue('Serial number')).doesNotExist('cert details not yet rendered');
    });

    test('it shows the certificate details after a successful fetch', async function (assert) {
      this.fetchStub.resolves(CERT_RESPONSE);
      await this.renderComponent();

      await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
      await fillIn(GENERAL.inputByAttr('minValidityValue'), '30');
      await click(GENERAL.submitButton);

      assert
        .dom(GENERAL.overviewCard.container('Retrieve cached certificate'))
        .doesNotExist('fetch form is gone after successful fetch');
      assert.true(this.fetchStub.calledOnce, 'API was called once');
      assert.dom(GENERAL.cardContainer('Certificate details')).exists();
      assert.dom(GENERAL.infoRowValue('Serial number')).exists().hasText(CERT_RESPONSE.serial_number);
    });

    test('it calls pkiExternalCaReadRoleCached with the correct role name and engine id', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
      await fillIn(GENERAL.inputByAttr('minValidityValue'), '60');
      await click(GENERAL.submitButton);

      const [roleName, engineId] = this.fetchStub.lastCall.args;
      assert.strictEqual(roleName, 'my-role', 'called with the role name from the model');
      assert.strictEqual(engineId, 'pki-external-ca', 'called with the engine id from the model');
    });

    test('it passes min_validity_duration in seconds when duration mode is used', async function (assert) {
      await this.renderComponent();

      await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
      await fillIn(GENERAL.inputByAttr('minValidityValue'), '2');
      await select(GENERAL.selectByAttr('unit'), 'h'); // 2 hours = 7200 seconds
      await click(GENERAL.submitButton);

      // Retrieve the initOverride callback the component passed and
      // invoke it so addQueryParams stub is called.
      const initOverride = this.fetchStub.lastCall.args[2];
      initOverride({});

      const [, payload] = this.addQueryParamsStub.lastCall.args;
      const expected = { identifiers: 'example.com', min_validity_duration: 7200 };
      assert.strictEqual(payload.min_validity_duration, 7200, '2h converted to 7200 seconds');
      assert.propEqual(payload, expected, 'callback has expected payload');
    });

    test('it passes min_validity_percentage when percentage mode is used', async function (assert) {
      await this.renderComponent();

      await click(GENERAL.radioByAttr('percentage'));
      await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
      await fillIn(GENERAL.inputByAttr('minValidityValue'), '50');
      await click(GENERAL.submitButton);

      // Retrieve the initOverride callback the component passed and
      // invoke it so addQueryParams stub is called.
      const initOverride = this.fetchStub.lastCall.args[2];
      initOverride({});

      const [, payload] = this.addQueryParamsStub.lastCall.args;
      const expected = { identifiers: 'example.com', min_validity_percentage: 50 };
      assert.strictEqual(payload.min_validity_percentage, 50, 'percentage payload is passed through');
      assert.propEqual(payload, expected, 'callback has expected payload');
    });

    // ---------------------------------------------------------------------------
    // Error handling — API failure surfaces an error on the form
    // ---------------------------------------------------------------------------

    test('it shows an error message when the API call fails', async function (assert) {
      const error = getErrorResponse({ errors: ['permission denied'] }, 403);
      this.fetchStub.rejects(error);

      await this.renderComponent();

      await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
      await fillIn(GENERAL.inputByAttr('minValidityValue'), '30');
      await click(GENERAL.submitButton);

      assert
        .dom(GENERAL.overviewCard.container('Retrieve cached certificate'))
        .exists('form is still shown after a failed fetch');
      assert.dom(GENERAL.messageError).exists().containsText('Error permission denied');
      assert.dom(GENERAL.cardContainer('Certificate details')).doesNotExist();
    });
  }
);
