/**
 * Copyright (c) HashiCorp, Inc.
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
    await click(ts.navLink('Secrets Sync'));
    await click(ts.cta.button);
    await click(ts.selectType('aws-sm'));
    await fillIn(ts.inputByAttr('name'), 'foo');
    await click(GENERAL.submitButton);
    assert.dom(ts.infoRowValue('Name')).hasText('foo', 'Destination details render after create success');

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

      await click(ts.navLink('Secrets Sync'));
      await click(ts.cta.button);
      await click(ts.selectType(type));

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
});
