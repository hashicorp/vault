/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import syncScenario from 'vault/mirage/scenarios/sync';
import syncHandlers from 'vault/mirage/handlers/sync';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { click, visit, fillIn, currentURL, currentRouteName } from '@ember/test-helpers';
import { PAGE as ts } from 'vault/tests/helpers/sync/sync-selectors';
import { syncDestinations } from 'vault/helpers/sync-destinations';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { DestinationType, CredentialType } from 'sync/utils/constants';
import sinon from 'sinon';

const SYNC_DESTINATIONS = syncDestinations();

// sync is an enterprise feature but since mirage is used the enterprise label has been intentionally omitted from the module name
module('Acceptance | sync | destinations (plural)', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.version = this.owner.lookup('service:version');
    this.version.features = ['Secrets Sync'];
    syncScenario(this.server);
    syncHandlers(this.server);
    return login();
  });

  test('it should create new destination', async function (assert) {
    // remove destinations from mirage so cta shows when 404 is returned
    this.server.db.syncDestinations.remove();

    this.server.get('/sys/activation-flags', () => {
      return {
        data: {
          activated: ['secrets-sync'],
          unactivated: [],
        },
      };
    });
    await click(GENERAL.navLink('Secrets'));
    await click(GENERAL.navLink('Secrets sync'));
    await click(ts.cta.button);
    await click(ts.selectType(DestinationType.AwsSm));
    await fillIn(GENERAL.inputByAttr('name'), 'foo');
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.infoRowValue('Name'))
      .hasText('foo', 'Destination details render after create success');

    await click(ts.breadcrumbLink('Destinations'));
    await click(ts.destinations.list.create);
    assert.strictEqual(
      currentURL(),
      '/vault/sync/secrets/destinations/create',
      'Toolbar action navigates to destinations create view'
    );
  });

  for (const destination of SYNC_DESTINATIONS) {
    const { type, defaultValues } = destination;
    test(`it should render default values for destination: ${type}`, async function (assert) {
      // remove destinations from mirage so cta shows when 404 is returned
      this.server.db.syncDestinations.remove();
      this.server.get('/sys/activation-flags', () => {
        return {
          data: {
            activated: ['secrets-sync'],
            unactivated: [],
          },
        };
      });

      await click(GENERAL.navLink('Secrets'));
      await click(GENERAL.navLink('Secrets sync'));
      await click(ts.cta.button);
      await click(ts.selectType(type));

      // Expand Advanced configuration accordion to access granularity field
      await click(GENERAL.accordionButton('Advanced configuration'));

      // check default values
      const attr = 'granularity';
      assert
        .dom(`${ts.inputGroupByAttr(attr)} input#${defaultValues[attr]}`)
        .isChecked(`${defaultValues[attr]} is checked`);
    });
  }

  test('it should filter destinations list', async function (assert) {
    await visit('vault/sync/secrets/destinations');
    assert.dom(GENERAL.listItemLink).exists({ count: 6 }, 'All destinations render');
    await click(`${ts.filter('type')} .ember-basic-dropdown-trigger`);
    await click(ts.searchSelect.option());
    assert.dom(GENERAL.listItemLink).exists({ count: 2 }, 'Destinations are filtered by type');
    await fillIn(ts.filter('name'), 'new');
    assert.dom(GENERAL.listItemLink).exists({ count: 1 }, 'Destinations are filtered by type and name');
    await click(ts.searchSelect.removeSelected);
    await fillIn(ts.filter('name'), 'gcp');
    assert.dom(GENERAL.listItemLink).exists({ count: 1 }, 'Destinations are filtered by name');
  });

  test('it should transition to correct routes when performing actions', async function (assert) {
    const routeName = (route) => `vault.cluster.sync.secrets.destinations.destination.${route}`;
    await visit('vault/sync/secrets/destinations');
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('details'));
    assert.strictEqual(currentRouteName(), routeName('details'), 'Navigates to details route');
    await click(ts.breadcrumbLink('Destinations'));
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('edit'));
    assert.strictEqual(currentRouteName(), routeName('edit'), 'Navigates to edit route');
  });

  // WIF (Workload Identity Federation) acceptance tests
  module('WIF credential type support', function (hooks) {
    hooks.beforeEach(function () {
      // Helper to clear destinations and activate sync feature
      this.clearDestinationsAndActivateSync = () => {
        this.server.db.syncDestinations.remove();
        this.server.get('/sys/activation-flags', () => {
          return {
            data: {
              activated: ['secrets-sync'],
              unactivated: [],
            },
          };
        });
      };

      // Helper to navigate to create destination form
      this.navigateToCreateDestination = async () => {
        await click(GENERAL.navLink('Secrets'));
        await click(GENERAL.navLink('Secrets sync'));
        await click(ts.cta.button);
      };
    });

    test('it should create AWS destination with WIF credentials', async function (assert) {
      this.clearDestinationsAndActivateSync();

      assert.expect(5);

      this.server.post('/sys/sync/destinations/aws-sm/:name', (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        assert.notOk('credential_type' in payload, 'credential_type not in payload');
        assert.notOk('access_key_id' in payload, 'account credentials not in payload');
        assert.ok('identity_token_audience' in payload, 'WIF credentials in payload');

        const { name } = req.params;
        const data = { ...payload, type: DestinationType.AwsSm, name };
        const record = schema.db.syncDestinations.insert(data);

        const { granularity, secret_name_template, custom_tags, ...connection_details } = record;
        return {
          data: {
            name: record.name,
            type: record.type,
            connection_details,
            options: {
              granularity_level: granularity,
              secret_name_template,
              custom_tags,
            },
          },
        };
      });

      await this.navigateToCreateDestination();
      await click(ts.selectType(DestinationType.AwsSm));

      // Switch to WIF
      await click(GENERAL.radioCardByAttr(CredentialType.WIF));

      // Fill in required fields
      await fillIn(GENERAL.inputByAttr('name'), 'wif-destination');
      await fillIn(GENERAL.inputByAttr('region'), 'us-west-1');
      await fillIn(GENERAL.inputByAttr('role_arn'), 'arn:aws:iam::123456789012:role/test-role');
      await fillIn(GENERAL.inputByAttr('identity_token_audience'), 'test-audience');

      await click(GENERAL.submitButton);

      assert.strictEqual(
        currentURL(),
        '/vault/sync/secrets/destinations/aws-sm/wif-destination/details',
        'Navigates to destination details after creation'
      );
      assert.dom(GENERAL.infoRowValue('Name')).hasText('wif-destination', 'Destination created successfully');
    });

    test('it should create Azure destination with WIF credentials', async function (assert) {
      this.clearDestinationsAndActivateSync();

      assert.expect(4);

      this.server.post('/sys/sync/destinations/azure-kv/:name', (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        assert.notOk('credential_type' in payload, 'credential_type not in payload');
        assert.notOk('client_secret' in payload, 'account credentials not in payload');
        assert.ok('identity_token_audience' in payload, 'WIF credentials in payload');

        const { name } = req.params;
        const data = { ...payload, type: DestinationType.AzureKv, name };
        const record = schema.db.syncDestinations.insert(data);

        const { granularity, secret_name_template, custom_tags, ...connection_details } = record;
        return {
          data: {
            name: record.name,
            type: record.type,
            connection_details,
            options: {
              granularity_level: granularity,
              secret_name_template,
              custom_tags,
            },
          },
        };
      });

      await this.navigateToCreateDestination();
      await click(ts.selectType(DestinationType.AzureKv));

      // Switch to WIF
      await click(GENERAL.radioCardByAttr(CredentialType.WIF));

      // Fill in required fields
      await fillIn(GENERAL.inputByAttr('name'), 'azure-wif');
      await fillIn(GENERAL.inputByAttr('key_vault_uri'), 'https://my-vault.vault.azure.net');
      await fillIn(GENERAL.inputByAttr('tenant_id'), 'tenant-id');
      await fillIn(GENERAL.inputByAttr('client_id'), 'client-id');
      await fillIn(GENERAL.inputByAttr('identity_token_audience'), 'test-audience');

      await click(GENERAL.submitButton);

      assert.strictEqual(
        currentURL(),
        '/vault/sync/secrets/destinations/azure-kv/azure-wif/details',
        'Navigates to destination details after creation'
      );
    });

    test('it should create GCP destination with WIF credentials', async function (assert) {
      this.clearDestinationsAndActivateSync();

      assert.expect(5);

      this.server.post('/sys/sync/destinations/gcp-sm/:name', (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        assert.notOk('credential_type' in payload, 'credential_type not in payload');
        assert.notOk('credentials' in payload, 'account credentials not in payload');
        assert.ok('identity_token_audience' in payload, 'WIF credentials in payload');
        assert.ok('service_account_email' in payload, 'service_account_email in payload');

        const { name } = req.params;
        const data = { ...payload, type: DestinationType.GcpSm, name };
        const record = schema.db.syncDestinations.insert(data);

        const { granularity, secret_name_template, custom_tags, ...connection_details } = record;
        return {
          data: {
            name: record.name,
            type: record.type,
            connection_details,
            options: {
              granularity_level: granularity,
              secret_name_template,
              custom_tags,
            },
          },
        };
      });

      await this.navigateToCreateDestination();
      await click(ts.selectType(DestinationType.GcpSm));

      // Switch to WIF
      await click(GENERAL.radioCardByAttr(CredentialType.WIF));

      // Fill in required fields
      await fillIn(GENERAL.inputByAttr('name'), 'gcp-wif');
      await fillIn(GENERAL.inputByAttr('project_id'), 'my-project');
      await fillIn(GENERAL.inputByAttr('service_account_email'), 'test@project.iam.gserviceaccount.com');
      await fillIn(GENERAL.inputByAttr('identity_token_audience'), 'test-audience');

      await click(GENERAL.submitButton);

      assert.strictEqual(
        currentURL(),
        '/vault/sync/secrets/destinations/gcp-sm/gcp-wif/details',
        'Navigates to destination details after creation'
      );
    });

    test('it should show credential type radio cards are checked by default for account', async function (assert) {
      this.clearDestinationsAndActivateSync();

      await this.navigateToCreateDestination();
      await click(ts.selectType(DestinationType.AwsSm));

      assert
        .dom(GENERAL.radioCardByAttr(CredentialType.ACCOUNT))
        .isChecked('Account credential type is selected by default');
      assert.dom(GENERAL.fieldByAttr('access_key_id')).exists('IAM credentials fields are visible');
    });

    test('it should display success message after creating WIF destination', async function (assert) {
      assert.expect(2);
      const flash = this.owner.lookup('service:flash-messages');
      const flashSuccessSpy = sinon.spy(flash, 'success');

      this.clearDestinationsAndActivateSync();

      this.server.post('/sys/sync/destinations/aws-sm/:name', (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        const { name } = req.params;
        const data = { ...payload, type: DestinationType.AwsSm, name };
        const record = schema.db.syncDestinations.insert(data);

        const { granularity, secret_name_template, custom_tags, ...connection_details } = record;
        return {
          data: {
            name: record.name,
            type: record.type,
            connection_details,
            options: {
              granularity_level: granularity,
              secret_name_template,
              custom_tags,
            },
          },
        };
      });

      await this.navigateToCreateDestination();
      await click(ts.selectType(DestinationType.AwsSm));

      // Switch to WIF
      await click(GENERAL.radioCardByAttr(CredentialType.WIF));

      // Fill in required fields
      await fillIn(GENERAL.inputByAttr('name'), 'test-wif');
      await fillIn(GENERAL.inputByAttr('region'), 'us-west-1');
      await fillIn(GENERAL.inputByAttr('role_arn'), 'arn:aws:iam::123456789012:role/test-role');
      await fillIn(GENERAL.inputByAttr('identity_token_audience'), 'test-audience');

      await click(GENERAL.submitButton);

      assert.true(flashSuccessSpy.calledOnce, 'Flash success message is called');
      const [flashMessage] = flashSuccessSpy.lastCall.args;
      assert.strictEqual(flashMessage, 'You have successfully created a sync destination.');
    });
  });
});
