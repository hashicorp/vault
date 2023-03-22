/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { create } from 'ember-cli-page-object';
import configPki from 'vault/tests/pages/components/pki/config-pki';

const component = create(configPki);

module('Integration | Component | config pki', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(async function () {
    this.owner.lookup('service:flash-messages').registerTypes(['success']);
    this.store = this.owner.lookup('service:store');
    this.config = await this.store.createRecord('pki/pki-config');
    this.mockConfigSave = function (saveFn) {
      const { tidyAttrs, crlAttrs, urlsAttrs } = this.config;
      return {
        save: saveFn,
        rollbackAttributes: () => {},
        tidyAttrs,
        crlAttrs,
        urlsAttrs,
        set: () => {},
      };
    };
  });

  const setupAndRender = async function (context, config, section = 'tidy') {
    context.set('config', config);
    context.set('section', section);
    await context.render(hbs`<Pki::ConfigPki @section={{this.section}} @config={{this.config}} />`);
  };

  test('it renders tidy section', async function (assert) {
    await setupAndRender(this, this.config);
    assert.ok(component.text.startsWith('You can tidy up the backend'));
    assert.notOk(component.hasTitle, 'No title for tidy section');
    assert.strictEqual(component.fields.length, 3, 'renders all three tidy fields');
    assert.ok(component.fields.objectAt(0).labelText, 'Tidy the Certificate Store');
    assert.ok(component.fields.objectAt(1).labelText, 'Tidy the Revocation List (CRL)');
    assert.ok(component.fields.objectAt(1).labelText, 'Safety buffer');
  });

  test('it renders crl section', async function (assert) {
    await setupAndRender(this, this.config, 'crl');
    assert.false(this.config.disable, 'CRL config defaults disable=false');
    assert.ok(component.hasTitle, 'renders form title');
    assert.strictEqual(component.title, 'Certificate Revocation List (CRL) config');
    assert.ok(
      component.text.startsWith('Set the duration for which the generated CRL'),
      'renders form subtext'
    );
    assert
      .dom('[data-test-ttl-form-label="CRL building enabled"]')
      .hasText('CRL building enabled', 'renders enabled field title');
    assert
      .dom('[data-test-ttl-form-subtext]')
      .hasText('The CRL will expire after', 'renders enabled field subtext');
    assert.dom('[data-test-input="expiry"] input').isChecked('defaults to enabling CRL build');
    assert.dom('[data-test-ttl-value="CRL building enabled"]').hasValue('3', 'default value is 3 (72h)');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('d', 'default unit value is days');
    await click('[data-test-input="expiry"] input');
    assert
      .dom('[data-test-ttl-form-subtext]')
      .hasText('The CRL will not be built.', 'renders disabled text when toggled off');

    // assert 'disable' attr on pki-config model updates with toggle
    assert.true(this.config.disable, 'when toggled off, sets CRL config to disable=true');
    await click('[data-test-input="expiry"] input');
    assert
      .dom('[data-test-ttl-form-subtext]')
      .hasText('The CRL will expire after', 'toggles back to enabled text');
    assert.false(this.config.disable, 'CRL config toggles back to disable=false');
  });

  test('it renders urls section', async function (assert) {
    await setupAndRender(this, this.config, 'urls');
    assert.notOk(component.hasTitle, 'No title for urls section');
    assert.strictEqual(component.fields.length, 3);
    assert.ok(component.fields.objectAt(0).labelText, 'Issuing certificates');
    assert.ok(component.fields.objectAt(1).labelText, 'CRL Distribution Points');
    assert.ok(component.fields.objectAt(2).labelText, 'OCSP Servers');
  });

  test('it calls save with the correct arguments for tidy', async function (assert) {
    assert.expect(3);
    const section = 'tidy';
    this.set('onRefresh', () => {
      assert.ok(true, 'refresh called');
    });
    this.set(
      'config',
      this.mockConfigSave((options) => {
        assert.strictEqual(options.adapterOptions.method, section, 'method passed to save');
        assert.deepEqual(
          options.adapterOptions.fields,
          ['tidyCertStore', 'tidyRevocationList', 'safetyBuffer'],
          'tidy fields passed to save'
        );
        return resolve();
      })
    );
    this.set('section', section);
    await render(
      hbs`<Pki::ConfigPki @section={{this.section}} @config={{this.config}} @onRefresh={{this.onRefresh}} />`
    );

    component.submit();
  });

  test('it calls save with the correct arguments for crl', async function (assert) {
    assert.expect(3);
    const section = 'crl';
    this.set('onRefresh', () => {
      assert.ok(true, 'refresh called');
    });
    this.set(
      'config',
      this.mockConfigSave((options) => {
        assert.strictEqual(options.adapterOptions.method, section, 'method passed to save');
        assert.deepEqual(options.adapterOptions.fields, ['expiry', 'disable'], 'CRL fields passed to save');
        return resolve();
      })
    );
    this.set('section', section);
    await render(
      hbs`<Pki::ConfigPki @section={{this.section}} @config={{this.config}} @onRefresh={{this.onRefresh}} />`
    );
    component.submit();
  });

  test('it correctly sets toggle when initial CRL config is disable=true', async function (assert) {
    assert.expect(3);
    // change default config attrs
    const configDisabled = this.config;
    configDisabled.expiry = '1m';
    configDisabled.disable = true;
    await setupAndRender(this, configDisabled, 'crl');
    assert.dom('[data-test-input="expiry"] input').isNotChecked('toggle disabled when CRL config disabled');
    await click('[data-test-input="expiry"] input');
    assert
      .dom('[data-test-ttl-value="CRL building enabled"]')
      .hasValue('1', 'when toggled on shows last set expired value');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('m', 'when toggled back on shows last set unit');
  });
});
