/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable ember/no-settled-after-test-helper */
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import syncScenario from 'vault/mirage/scenarios/sync';
import syncHandlers from 'vault/mirage/handlers/sync';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { Response } from 'miragejs';
import { dateFormat } from 'core/helpers/date-format';

const { title, tab, overviewCard, cta, overview, pagination, emptyStateTitle, emptyStateMessage } = PAGE;

module('Integration | Component | sync | Page::Overview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';
    this.version.features = ['Secrets Sync'];
    syncScenario(this.server);
    syncHandlers(this.server);

    const store = this.owner.lookup('service:store');
    this.destinations = await store.query('sync/destination', {});
    this.activatedFeatures = ['secrets-sync'];

    this.renderComponent = () =>
      render(
        hbs`<Secrets::Page::Overview @destinations={{this.destinations}} @totalVaultSecrets={{7}} @activatedFeatures={{this.activatedFeatures}} @adapterError={{null}} />`,
        {
          owner: this.engine,
        }
      );
  });

  test('it should render header, tabs and toolbar for overview state', async function (assert) {
    await this.renderComponent();

    assert.dom(title).hasText('Secrets Sync', 'Page title renders');
    assert.dom(cta.button).doesNotExist('CTA does not render');
    assert.dom(tab('Overview')).hasText('Overview', 'Overview tab renders');
    assert.dom(tab('Destinations')).hasText('Destinations', 'Destinations tab renders');
    assert.dom(overview.createDestination).hasText('Create new destination', 'Toolbar action renders');
  });

  module('community', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'community';
      this.version.features = [];
      this.destinations = [];
      this.activatedFeatures = [];
    });

    test('it should show an upsell CTA', async function (assert) {
      await this.renderComponent();

      assert
        .dom(title)
        .hasText('Secrets Sync Enterprise feature', 'page title indicates feature is only for Enterprise');
      assert.dom(cta.button).doesNotExist();
    });
  });

  module('ent', function (hooks) {
    hooks.beforeEach(function () {
      this.destinations = [];
    });

    test('it should show an upsell CTA if license does NOT have the secrets sync feature', async function (assert) {
      this.version.features = [];
      await this.renderComponent();

      assert
        .dom(title)
        .hasText('Secrets Sync Premium feature', 'title indicates feature is only for Premium');
      assert.dom(cta.button).doesNotExist();
      assert.dom(cta.summary).exists();
    });

    test('it should show create CTA if license has the secrets sync feature', async function (assert) {
      this.version.features = ['Secrets Sync'];
      await this.renderComponent();

      assert.dom(title).hasText('Secrets Sync');
      assert.dom(cta.button).hasText('Create first destination', 'CTA action renders');
      assert.dom(cta.summary).exists();
    });
  });

  module('secrets sync not activated', function (hooks) {
    hooks.beforeEach(async function () {
      this.activatedFeatures = [];
    });

    test('it should show the opt-in banner', async function (assert) {
      await this.renderComponent();

      assert.dom(overview.optInBanner).exists('Opt-in banner is shown');
    });

    test('it should navigate to the opt-in modal', async function (assert) {
      await this.renderComponent();

      await click(overview.optInBannerEnable);

      assert.dom(overview.optInModal).exists('Opt-in modal is shown');
      assert.dom(overview.optInConfirm).isDisabled('Confirm button is disabled when checkbox is unchecked');

      await click(overview.optInCheck);
      assert.dom(overview.optInConfirm).isNotDisabled('confirm button is enabled once checkbox is checked');
    });

    test('it should make a POST to activate the feature', async function (assert) {
      assert.expect(1);

      await this.renderComponent();

      this.server.post('/sys/activation-flags/secrets-sync/activate', () => {
        assert.true(true, 'POST to secrets-sync/activate is called');
        return {};
      });

      await this.renderComponent();

      await click(overview.optInBannerEnable);
      await click(overview.optInCheck);
      await click(overview.optInConfirm);
    });

    test('it shows an error if activation fails', async function (assert) {
      await this.renderComponent();

      this.server.post('/sys/activation-flags/secrets-sync/activate', () => new Response(403));

      await click(overview.optInBannerEnable);
      await click(overview.optInCheck);
      await click(overview.optInConfirm);

      assert.dom(overview.optInError).exists('shows an error banner');
      assert.dom(overview.optInBanner).exists('banner is visible so user can try to opt-in again');
    });
  });

  module('secrets sync activated', function () {
    test('it should hide the opt-in banner', async function (assert) {
      await this.renderComponent();

      assert.dom(overview.optInBanner).doesNotExist();
    });
  });

  module('with no destinations', function (hooks) {
    hooks.beforeEach(function () {
      this.destinations = [];
    });

    test('it should show a CTA', async function (assert) {
      await this.renderComponent();

      assert.dom(cta.button).exists();

      assert
        .dom(overview.createDestination)
        .doesNotExist('create new destination link is hidden if there are no pre-existing destinations');
    });

    test('it should hide the overview cards', async function (assert) {
      await this.renderComponent();

      assert
        .dom(overviewCard.title('Secrets by destination'))
        .doesNotExist('it does not show secrets by destination card if there are no destinations');
      assert
        .dom(overviewCard.title('Total destinations'))
        .doesNotExist('it does not show total destinations card if there are no destinations');
      assert
        .dom(overviewCard.title('Total secrets'))
        .doesNotExist('it does not show total secrets card if there are no destinations');
    });
  });

  module('with destinations', function () {
    test('it should show the overview cards', async function (assert) {
      await this.renderComponent();

      assert.dom(overviewCard.title('Secrets by destination')).exists();
      assert.dom(overviewCard.title('Total destinations')).exists();
      assert.dom(overviewCard.title('Total secrets')).exists();
    });

    module('Secrets by destination table', function () {
      test('it should show the table with correct columns, data, badges, and actions', async function (assert) {
        const { icon, name, badge, total, updated, actionToggle, action } = overview.table;
        const updatedDate = dateFormat(
          [new Date('2023-09-20T10:51:53.961861096-04:00'), 'MMMM do yyyy, h:mm:ss a'],
          {}
        );

        await this.renderComponent();

        assert
          .dom(overviewCard.title('Secrets by destination'))
          .hasText('Secrets by destination', 'Overview card title renders for table');
        assert.dom(icon(0)).hasAttribute('data-test-icon', 'aws-color', 'Destination icon renders');
        assert.dom(name(0)).hasText('destination-aws', 'Destination name renders');

        assert.dom(badge(0)).hasText('All synced', 'All synced badge renders');
        assert
          .dom(badge(0))
          .hasClass('hds-badge--color-success', 'Correct color renders for All synced badge');

        assert
          .dom(badge(1))
          .doesNotExist('Status badge does not render for destination with no associations');

        assert.dom(badge(2)).hasText('1 Unsynced', 'Unsynced badge renders');
        assert.dom(badge(2)).hasClass('hds-badge--color-neutral', 'Correct color renders for unsynced badge');

        assert.dom(total(0)).hasText('1', '# of external secrets renders');
        assert.dom(updated(0)).hasText(updatedDate, 'Last updated datetime renders');

        assert
          .dom(total(1))
          .hasText('0', '# of external secrets renders for destination with no associations');
        assert
          .dom(updated(1))
          .hasText('—', 'Last updated placeholder renders for destination with no associations');

        await click(actionToggle(0));
        assert.dom(action('sync')).hasText('Sync secrets', 'Sync action renders');
        assert.dom(action('details')).hasText('View synced secrets', 'View synced secrets action renders');
      });

      test('it should paginate the table', async function (assert) {
        await this.renderComponent();

        const { name, row } = overview.table;
        assert.dom(row).exists({ count: 3 }, 'Correct number of table rows render based on page size');
        assert.dom(name(0)).hasText('destination-aws', 'First destination renders on page 1');

        await click(pagination.next);
        await settled();
        assert
          .dom(overview.table.row)
          .exists({ count: 3 }, 'New items are fetched and rendered on page change');
        assert.dom(name(0)).hasText('destination-gcp', 'First destination renders on page 2');
      });

      test('it should show an empty state if there is an error fetching associations', async function (assert) {
        this.server.get('/sys/sync/destinations/:type/:name/associations', () => {
          return new Response(403, {}, { errors: ['Permission denied'] });
        });

        await this.renderComponent();

        assert.dom(emptyStateTitle).hasText('Error fetching information', 'Empty state title renders');
        assert
          .dom(emptyStateMessage)
          .hasText(
            'Ensure that the policy has access to read sync associations.',
            'Empty state message renders'
          );
      });
    });

    test('it should show the Totals cards', async function (assert) {
      await this.renderComponent();

      const { title, description, action, content } = overviewCard;
      const cardData = [
        {
          cardTitle: 'Total destinations',
          subText: 'The total number of connected destinations',
          actionText: 'Create new',
          count: '6',
        },
        {
          cardTitle: 'Total secrets',
          subText:
            'The total number of secrets that have been synced from Vault. One secret will be counted as one sync client.',
          actionText: 'View billing',
          count: '7',
        },
      ];

      cardData.forEach(({ cardTitle, subText, actionText, count }) => {
        assert.dom(title(cardTitle)).hasText(cardTitle, `${cardTitle} card title renders`);
        assert.dom(description(cardTitle)).hasText(subText, ` ${cardTitle} card description renders`);
        assert.dom(content(cardTitle)).hasText(count, 'Total count renders');
        assert.dom(action(cardTitle)).hasText(actionText, 'Card action renders');
      });
    });
  });
});
