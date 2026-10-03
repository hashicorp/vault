/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import KmipConfigForm from 'vault/forms/secrets/kmip/config';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | kmip | Page::Configure', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.backend = 'kmip-test';
    this.owner.lookup('service:secret-mount-path').update(this.backend);

    const { secrets } = this.owner.lookup('service:api');
    this.apiStub = sinon.stub(secrets, 'kmipConfigure').resolves();
    this.flashStub = sinon.stub(this.owner.lookup('service:flashMessages'), 'success');
    this.routerStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    this.config = {
      default_tls_client_key_bits: 256,
      default_tls_client_key_type: 'ec',
      default_tls_client_ttl: '336h',
      listen_addrs: '127.0.0.1:5696',
      server_hostnames: ['localhost'],
      server_ips: ['0.0.0.0'],
      tls_ca_key_bits: 256,
      tls_ca_key_type: 'ec',
      tls_min_version: 'tls12',
    };
    this.form = new KmipConfigForm({}, { isNew: true });

    this.inputFor = (field, index = 0) => `${GENERAL.inputByAttr(field)} ${GENERAL.stringListByIdx(index)}`;

    this.renderComponent = () => render(hbs`<Page::Configure @form={{this.form}} />`, { owner: this.engine });
  });

  test('it should create new config', async function (assert) {
    await this.renderComponent();

    await fillIn(this.inputFor('server_hostnames'), 'localhost');
    await fillIn(this.inputFor('server_ips'), '0.0.0.0');
    await click(GENERAL.submitButton);

    assert.true(this.apiStub.calledWith(this.backend, this.config), 'API called with correct params');
    assert.true(
      this.flashStub.calledWith('Successfully configured KMIP engine'),
      'Success flash message shown'
    );
    assert.true(
      this.routerStub.calledWith('vault.cluster.secrets.backend.kmip.configuration'),
      'Transitions to configuration page'
    );
  });

  test('it should edit existing config', async function (assert) {
    this.form = new KmipConfigForm(this.config);

    await this.renderComponent();

    assert
      .dom(this.inputFor('server_hostnames'))
      .hasValue('localhost', 'server_hostnames is populated with config value');
    assert.dom(this.inputFor('server_ips')).hasValue('0.0.0.0', 'server_ips is populated with config value');

    await fillIn(this.inputFor('server_hostnames', 1), 'foobar');
    await fillIn(this.inputFor('server_ips', 1), '127.0.0.1');
    await click(GENERAL.submitButton);

    this.config.server_hostnames.push('foobar');
    this.config.server_ips.push('127.0.0.1');
    assert.true(this.apiStub.calledWith(this.backend, this.config), 'API called with correct params');
    assert.true(
      this.flashStub.calledWith('Successfully configured KMIP engine'),
      'Success flash message shown'
    );
    assert.true(
      this.routerStub.calledWith('vault.cluster.secrets.backend.kmip.configuration'),
      'Transitions to configuration page'
    );
  });

  test('it should handle errors', async function (assert) {
    this.apiStub.rejects(getErrorResponse({ errors: ['Invalid configuration provided'] }, 400));

    await this.renderComponent();
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.inlineError)
      .hasText('There was an error submitting this form.', 'Displays inline error message');
    assert
      .dom(GENERAL.messageDescription)
      .hasText('Invalid configuration provided', 'Displays error from API');
  });
});
