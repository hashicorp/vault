/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupDataStubs } from 'vault/tests/helpers/sync/setup-hooks';
import hbs from 'htmlbars-inline-precompile';
import { render, click, fillIn } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { selectChoose } from 'ember-power-select/test-support';
import { Response } from 'miragejs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const { destinations, searchSelect, messageError } = PAGE;
const { mountSelect, mountInput, successMessage } = destinations.sync;

module('Integration | Component | sync | Secrets::Page::Destinations::Destination::Sync', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);
  setupDataStubs(hooks);

  hooks.beforeEach(async function () {
    const api = this.owner.lookup('service:api');
    this.mountsStub = sinon.stub(api.sys, 'internalUiListEnabledVisibleMounts').resolves({
      secret: {
        'my-kv/': { type: 'kv', options: { version: '2' } },
        'my-db/': { type: 'database', options: {} },
        'transit/': { type: 'transit', options: {} }, // Should be filtered out
      },
    });

    this.secretsStub = sinon.stub(api.secrets, 'kvV2List').resolves({ keys: ['my-path/', 'my-secret'] });

    this.databaseStub = sinon.stub(api.secrets, 'databaseListStaticRoles').resolves({ keys: ['my-role'] });

    await render(
      hbs`<Secrets::Page::Destinations::Destination::Sync @destination={{this.destination}} @capabilities={{this.capabilities}} />`,
      {
        owner: this.engine,
      }
    );
  });

  test('it should fetch and render all supported mounts with type indicators', async function (assert) {
    assert.expect(3);
    await selectChoose(mountSelect, '.ember-power-select-option', 0);
    assert.dom(mountSelect).includesText('my-kv/', 'kv mount is fetched and renders in super select');
    assert.dom(mountSelect).includesText('KV v2', 'kv mount shows type indicator');
    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    assert.dom(mountSelect).includesText('my-db/', 'database mount is fetched and renders in super select');
  });

  test('it should auto-detect secret type from selected mount', async function (assert) {
    assert.expect(4);
    // Select KV mount
    await selectChoose(mountSelect, '.ember-power-select-option', 0);
    assert.dom(GENERAL.suggestion.input('kv')).exists('KV suggestion input renders for KV mount');
    assert
      .dom(GENERAL.suggestion.input('database'))
      .doesNotExist('Database suggestion input does not render');

    // Select Database mount
    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    assert
      .dom(GENERAL.suggestion.input('database'))
      .exists('Database suggestion input renders for Database mount');
    assert.dom(GENERAL.suggestion.input('kv')).doesNotExist('KV suggestion input does not render');
  });

  test('it should filter unsupported mount types', async function (assert) {
    // Transit mount should not appear in the dropdown
    await click(mountSelect);
    const options = this.element.querySelectorAll('.ember-power-select-option');
    const optionTexts = Array.from(options).map((el) => el.textContent.trim());
    assert.ok(
      optionTexts.some((text) => text.includes('my-kv/')),
      'KV mount is available'
    );
    assert.ok(
      optionTexts.some((text) => text.includes('my-db/')),
      'Database mount is available'
    );
    assert.notOk(
      optionTexts.some((text) => text.includes('transit')),
      'Transit mount is filtered out'
    );
    // Verify icons are rendered in dropdown options
    assert.dom('.ember-power-select-option:first-child .hds-icon').exists('Icons render in dropdown options');
  });

  test('it should render secret suggestions for selected mount', async function (assert) {
    assert
      .dom(GENERAL.suggestion.input('kv'))
      .isDisabled('Secret input is disabled until a mount is selected');
    await selectChoose(mountSelect, '.ember-power-select-option', 0);
    await click(GENERAL.suggestion.input('kv'));
    assert.dom(searchSelect.option()).hasText('my-path/', 'Nested secret path renders');
    assert.dom(searchSelect.option(1)).hasText('my-secret', 'Secret renders');
  });

  test('it should render secret suggestions for nested paths', async function (assert) {
    await selectChoose(mountSelect, '.ember-power-select-option', 0);
    this.secretsStub.resolves({ keys: ['nested-secret'] });
    await click(GENERAL.suggestion.input('kv'));
    await click(searchSelect.option());
    assert
      .dom(searchSelect.option())
      .hasText('nested-secret', 'Suggestions render for secret at nested path');
  });

  test('it should sync secret', async function (assert) {
    assert.expect(8);

    const { type, name } = this.destination;
    this.server.post(`/sys/sync/destinations/${type}/${name}/associations/set`, (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = { mount: 'my-kv', secret_name: 'my-secret' };
      assert.deepEqual(data, expected, 'Sync request made with mount and secret name');
      return { data: { associated_secrets: { 'my-kv_12345': data } } };
    });
    assert.dom(GENERAL.submitButton).isDisabled('Submit button is disabled when mount is not selected');
    assert.dom(GENERAL.cancelButton).hasText('Cancel', 'Cancel button renders');
    await selectChoose(mountSelect, '.ember-power-select-option', 0);
    assert.dom(GENERAL.submitButton).isDisabled('Submit button is disabled when secret is not selected');
    await click(GENERAL.suggestion.input('kv'));
    await click(searchSelect.option(1));
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.cancelButton).hasText('View synced secrets', 'view secrets tertiary renders');
    assert.dom(GENERAL.suggestion.input('kv')).hasNoValue('Secret path is unset after submit success');
    assert.dom(GENERAL.submitButton).isDisabled('Submit button is disabled');
    assert
      .dom(successMessage)
      .includesText('Sync operation successfully initiated for my-secret.', 'Success banner renders');
  });

  test('it should allow manual mount path input if kv mounts are not returned', async function (assert) {
    assert.expect(1);

    this.mountsStub.resolves({ secret: { 'cubbyhole/': { type: 'cubbyhole' } } });

    const { type, name } = this.destination;
    this.server.post(`/sys/sync/destinations/${type}/${name}/associations/set`, (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = { mount: 'my-kv', secret_name: 'my-secret' };
      assert.deepEqual(data, expected, 'Sync request made with mount and secret name');
      return {};
    });

    await render(hbs`<Secrets::Page::Destinations::Destination::Sync @destination={{this.destination}} />`, {
      owner: this.engine,
    });

    // Manually enter mount path since no KV mounts are available
    await fillIn(mountInput, 'my-kv');
    await click(GENERAL.suggestion.input('kv'));
    await click(searchSelect.option(1));
    await click(GENERAL.submitButton);
  });

  test('it should render alert banner on sync error', async function (assert) {
    assert.expect(1);

    const { type, name } = this.destination;
    const error = 'Secret not found. Please provide full path to existing secret';
    this.server.post(`/sys/sync/destinations/${type}/${name}/associations/set`, () => {
      return new Response(400, {}, { errors: [error] });
    });

    await selectChoose(mountSelect, '.ember-power-select-option', 0);
    await click(GENERAL.suggestion.input('kv'));
    await click(searchSelect.option(1));
    await click(GENERAL.submitButton);
    assert.dom(messageError).hasTextContaining(error, 'Error renders in alert banner');
  });

  test('it should sync database role', async function (assert) {
    assert.expect(3);

    const { type, name } = this.destination;
    this.server.post(`/sys/sync/destinations/${type}/${name}/associations/set`, (schema, req) => {
      const data = JSON.parse(req.requestBody);
      const expected = { mount: 'my-db', secret_name: 'static-roles/my-role' };
      assert.deepEqual(data, expected, 'Sync request made with mount and role name');
      return { data: { associated_secrets: { 'my-db_12345': data } } };
    });

    await selectChoose(mountSelect, '.ember-power-select-option', 1);
    await click(GENERAL.suggestion.input('database'));
    await click(searchSelect.option());
    await click(GENERAL.submitButton);
    assert
      .dom(successMessage)
      .includesText(
        'Sync operation successfully initiated for static-roles/my-role',
        'Success banner renders for database'
      );
    assert
      .dom(`${successMessage} a`)
      .exists('External link renders for database type (supportsExternalLink: true)');
  });
});
