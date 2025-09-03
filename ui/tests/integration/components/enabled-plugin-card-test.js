/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, triggerKeyEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | enabled-plugin-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.setMountType = sinon.spy();
    this.type = {
      type: 'kv',
      displayName: 'KV',
      glyph: 'key-values',
    };
    this.version = this.owner.lookup('service:version');
  });

  test('it renders basic plugin card', async function (assert) {
    await render(hbs`
      <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
    `);

    assert.dom(GENERAL.cardContainer('kv')).exists('renders plugin card');
    assert.dom(`${GENERAL.cardContainer('kv')} h3`).hasText('KV', 'displays plugin name');
    assert.dom(`${GENERAL.cardContainer('kv')} [data-test-icon]`).exists('displays plugin icon');
  });

  test('it handles click to select plugin', async function (assert) {
    await render(hbs`
      <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
    `);

    await click(GENERAL.cardContainer('kv'));

    assert.ok(this.setMountType.calledOnce, 'calls setMountType on click');
    assert.ok(this.setMountType.calledWith('kv'), 'passes correct plugin type');
  });

  test('it handles keyboard navigation', async function (assert) {
    await render(hbs`
      <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
    `);

    await triggerKeyEvent(GENERAL.cardContainer('kv'), 'keydown', 'Enter');
    assert.ok(this.setMountType.calledOnce, 'calls setMountType on Enter key');

    this.setMountType.resetHistory();
    await triggerKeyEvent(GENERAL.cardContainer('kv'), 'keydown', ' ');
    assert.ok(this.setMountType.calledOnce, 'calls setMountType on Space key');

    this.setMountType.resetHistory();
    await triggerKeyEvent(GENERAL.cardContainer('kv'), 'keydown', 'Escape');
    assert.notOk(this.setMountType.called, 'does not call setMountType on other keys');
  });

  module('enterprise features', function () {
    test('it shows enterprise badge when required and not enterprise', async function (assert) {
      this.type = {
        ...this.type,
        requiresEnterprise: true,
      };
      sinon.stub(this.version, 'isEnterprise').value(false);

      await render(hbs`
        <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
      `);

      assert.dom(`${GENERAL.cardContainer('kv')} .hds-badge`).hasText('Enterprise', 'shows enterprise badge');
      assert.dom(GENERAL.cardContainer('kv')).hasClass('disabled', 'adds disabled class');
    });

    test('it does not show enterprise badge when enterprise license present', async function (assert) {
      this.type = {
        ...this.type,
        requiresEnterprise: true,
      };
      sinon.stub(this.version, 'isEnterprise').value(true);

      await render(hbs`
        <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
      `);

      assert.dom(`${GENERAL.cardContainer('kv')} .hds-badge`).doesNotExist('does not show enterprise badge');
      assert.dom(GENERAL.cardContainer('kv')).doesNotHaveClass('disabled', 'does not add disabled class');
    });

    test('it shows enterprise badge for required feature not available', async function (assert) {
      this.type = {
        ...this.type,
        requiredFeature: 'Advanced Data Protection',
      };
      this.version.set('features', []);

      await render(hbs`
        <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
      `);

      assert.dom(`${GENERAL.cardContainer('kv')} .hds-badge`).hasText('Enterprise', 'shows enterprise badge');
      assert.dom(GENERAL.cardContainer('kv')).hasClass('disabled', 'adds disabled class');
    });

    test('it does not show enterprise badge when feature is available', async function (assert) {
      this.type = {
        ...this.type,
        requiredFeature: 'Advanced Data Protection',
      };
      this.version.set('features', ['Advanced Data Protection']);

      await render(hbs`
        <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
      `);

      assert.dom(`${GENERAL.cardContainer('kv')} .hds-badge`).doesNotExist('does not show enterprise badge');
      assert.dom(GENERAL.cardContainer('kv')).doesNotHaveClass('disabled', 'does not add disabled class');
    });

    test('it does not call setMountType when disabled', async function (assert) {
      this.type = {
        ...this.type,
        requiresEnterprise: true,
      };
      sinon.stub(this.version, 'isEnterprise').value(false);

      await render(hbs`
        <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
      `);

      await click(GENERAL.cardContainer('kv'));
      assert.notOk(this.setMountType.called, 'does not call setMountType when disabled');
    });
  });

  module('deprecation status', function () {
    test('it shows deprecation badge when deprecated', async function (assert) {
      this.type = {
        ...this.type,
        deprecationStatus: 'deprecated',
      };

      await render(hbs`
        <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
      `);

      assert
        .dom(`${GENERAL.cardContainer('kv')} .hds-badge`)
        .hasText('Deprecated', 'shows deprecation badge');
    });

    test('it shows pending removal badge', async function (assert) {
      this.type = {
        ...this.type,
        deprecationStatus: 'pending removal',
      };

      await render(hbs`
        <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
      `);

      assert
        .dom(`${GENERAL.cardContainer('kv')} .hds-badge`)
        .hasText('Pending removal', 'shows pending removal badge');
    });

    test('it does not show badge when supported', async function (assert) {
      this.type = {
        ...this.type,
        deprecationStatus: 'supported',
      };

      await render(hbs`
        <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
      `);

      assert
        .dom(`${GENERAL.cardContainer('kv')} .hds-badge`)
        .doesNotExist('does not show badge for supported');
    });

    test('it does not show badge when no deprecation status', async function (assert) {
      await render(hbs`
        <EnabledPluginCard @type={{this.type}} @setMountType={{this.setMountType}} />
      `);

      assert
        .dom(`${GENERAL.cardContainer('kv')} .hds-badge`)
        .doesNotExist('does not show badge when no status');
    });
  });
});
