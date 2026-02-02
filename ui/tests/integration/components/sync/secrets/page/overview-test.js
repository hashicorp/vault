/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable ember/no-settled-after-test-helper */
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, settled, findAll } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import syncScenario from 'vault/mirage/scenarios/sync';
import syncHandlers from 'vault/mirage/handlers/sync';
import sinon from 'sinon';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { Response } from 'miragejs';
import { dateFormat } from 'core/helpers/date-format';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { listDestinationsTransform } from 'sync/utils/api-transforms';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const { tab, overviewCard, cta, overview, emptyStateTitle, emptyStateMessage } = PAGE;

module('Integration | Component | sync | Page::Overview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    // allow capabilities as root by default to allow users to POST to the secrets-sync/activate endpoint
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.version = this.owner.lookup('service:version');
    this.api = this.owner.lookup('service:api');
    this.flags = this.owner.lookup('service:flags');

    syncScenario(this.server);
    syncHandlers(this.server);

    const destinations = await this.api.sys.systemListSyncDestinations(true);
    this.destinations = listDestinationsTransform(destinations);

    this.setup = ({
      canActivate = false,
      isActivated = false,
      isHvdManaged = false,
      isEnterprise = false,
      hasFeature = false,
      hasDestinations = false,
    } = {}) => {
      // permissions are checked in the route model and passed to the page component
      this.canActivate = canActivate;

      // sync has to be activated to use the feature
      this.flags.activatedFlags = isActivated ? ['secrets-sync'] : [];

      // cluster is HVD managed
      this.flags.featureFlags = isHvdManaged ? ['VAULT_CLOUD_ADMIN_NAMESPACE'] : [];

      // cluster is an enterprise version, HVD managed clusters are also enterprise
      this.version.type = isEnterprise || isHvdManaged ? 'enterprise' : 'community';

      // self-managed enterprise clusters need to have the sync feature on the license
      // (does not apply to HVD managed clusters)
      this.version.features = isEnterprise && hasFeature ? ['Secrets Sync'] : [];

      if (!hasDestinations) this.destinations = [];
    };

    this.renderComponent = () => {
      return render(
        hbs`<Secrets::Page::Overview
        @destinations={{this.destinations}}
        @totalVaultSecrets={{7}}
        @canActivateSecretsSync={{this.canActivate}} />`,
        { owner: this.engine }
      );
    };
  });

  // navigating to this route is hidden for CE, but if for some reason this template were to show...
  test('it should hide opt-in banner for community', async function (assert) {
    this.setup();
    await this.renderComponent();
    assert.dom(overview.optInBanner.container).doesNotExist();
    assert.dom(cta.summary).exists();
  });

  test('it should render header, tabs and toolbar for overview state if destinations exist', async function (assert) {
    this.setup({
      isActivated: true,
      isEnterprise: true,
      hasFeature: true,
      hasDestinations: true,
    });
    await this.renderComponent();

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Secrets Sync', 'Page title renders');
    assert.dom(cta.summary).doesNotExist('CTA does not render');
    assert.dom(tab('Overview')).hasText('Overview', 'Overview tab renders');
    assert.dom(tab('Destinations')).hasText('Destinations', 'Destinations tab renders');
    assert.dom(overview.createDestination).hasText('Create new destination', 'Toolbar action renders');
  });

  // HVD MANAGED CLUSTERS
  test('it should show the opt-in banner if feature is not activated', async function (assert) {
    this.setup({ isHvdManaged: true });
    await this.renderComponent();

    assert.dom(overview.optInBanner.container).exists('Opt-in banner is shown');
  });

  test('it should show CTA if feature is activated', async function (assert) {
    this.setup({ isActivated: true, isHvdManaged: true });
    await this.renderComponent();

    assert.dom(overview.optInBanner.container).doesNotExist('Opt-in banner is not shown');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Secrets Sync');
    assert.dom(GENERAL.badge('Plus feature')).hasText('Plus feature', 'Plus feature badge renders');
    assert.dom(cta.button).hasText('Create first destination', 'CTA action renders');
    assert.dom(cta.summary).exists();
  });

  test('it should show activation error if cluster is not Plus tier', async function (assert) {
    this.setup({ canActivate: true, isHvdManaged: true });
    await this.renderComponent();

    this.server.post(
      '/sys/activation-flags/secrets-sync/activate',
      () => new Response(403, {}, { errors: ['Something bad happened'] })
    );

    await click(overview.optInBanner.enable);
    await click(overview.activationModal.checkbox);
    await click(overview.activationModal.confirm);

    assert.dom(overview.optInError).exists({ count: 2 }, 'shows the API and custom tier error banners');

    const errorBanners = findAll(overview.optInError);

    assert.dom(errorBanners[0]).containsText('Something bad happened', 'shows the API error message');

    assert
      .dom(errorBanners[1])
      .containsText(
        'Error Secrets Sync is available for Plus tier clusters only. Please check the tier of your cluster to enable Secrets Sync.',
        'shows the custom tier-related error message'
      );

    assert.dom(overview.optInBanner.container).exists('banner is visible so user can try to opt-in again');
  });

  // ENTERPRISE CLUSTERS
  test('it should show create CTA if activated and license has the secrets sync feature', async function (assert) {
    this.setup({ isActivated: true, isEnterprise: true, hasFeature: true });

    this.version.features = ['Secrets Sync'];
    this.isActivated = true;
    await this.renderComponent();

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Secrets Sync');
    assert.dom(cta.button).hasText('Create first destination', 'CTA action renders');
    assert.dom(cta.summary).exists();
  });

  test('it should show the opt-in banner without permissions to activate', async function (assert) {
    this.setup({ isEnterprise: true, hasFeature: true });
    await this.renderComponent();

    assert
      .dom(overview.optInBanner.description)
      .hasText(
        'To use this feature, specific activation is required. Please contact your administrator to activate.'
      );
    assert.dom(overview.optInBanner.enable).doesNotExist('Opt-in enable button does not show');
    assert.dom(overview.optInBanner.dismiss).doesNotExist('dismiss opt-in banner does not show');
  });

  test('it should show the opt-in banner with activate description with permissions to activate', async function (assert) {
    this.setup({ canActivate: true, isEnterprise: true, hasFeature: true });
    await this.renderComponent();

    assert.dom(overview.optInBanner.container).exists('Opt-in banner is shown');
    assert
      .dom(overview.optInBanner.description)
      .hasText(
        "To use this feature, specific activation is required. Please review the feature documentation and enable it. If you're upgrading from beta, your previous data will be accessible after activation."
      );
    assert.dom(overview.optInBanner.dismiss).exists('dismiss opt-in banner shows');
  });

  test('it should navigate to the opt-in modal', async function (assert) {
    this.setup({ canActivate: true, isEnterprise: true, hasFeature: true });
    await this.renderComponent();

    await click(overview.optInBanner.enable);

    assert.dom(overview.activationModal.container).exists('Opt-in modal is shown');
  });

  test('it shows an error if activation fails', async function (assert) {
    this.setup({ canActivate: true, isEnterprise: true, hasFeature: true });
    await this.renderComponent();

    this.server.post(
      '/sys/activation-flags/secrets-sync/activate',
      () => new Response(403, {}, { errors: ['Something bad happened'] })
    );

    await click(overview.optInBanner.enable);
    await click(overview.activationModal.checkbox);
    await click(overview.activationModal.confirm);

    assert
      .dom(overview.optInError)
      .exists({ count: 1 })
      .containsText('Something bad happened', 'shows an error banner with error message from the API');
    assert.dom(overview.optInBanner.container).exists('banner is visible so user can try to opt-in again');
  });

  test('it should clear activation errors when the user tries to opt-in again', async function (assert) {
    this.setup({ canActivate: true, isEnterprise: true, hasFeature: true });
    // don't worry about transitioning the route in this test
    sinon.stub(this.owner.lookup('service:router'), 'refresh');

    await this.renderComponent();

    let callCount = 0;

    // first call fails, second call succeeds
    this.server.post('/sys/activation-flags/secrets-sync/activate', () => {
      callCount++;
      if (callCount === 1) {
        return new Response(403, {}, { errors: ['Something bad happened'] });
      } else {
        return {};
      }
    });

    await click(overview.optInBanner.enable);
    await click(overview.activationModal.checkbox);
    await click(overview.activationModal.confirm);

    assert.dom(overview.optInError).exists('shows an error banner');

    await click(overview.optInBanner.enable);
    await click(overview.activationModal.checkbox);
    await click(overview.activationModal.confirm);

    assert.dom(overview.optInError).doesNotExist('error banner is cleared upon trying to opt-in again');
  });

  test('it should show a CTA if activated and no destinations', async function (assert) {
    this.setup({ isActivated: true, isEnterprise: true, hasFeature: true });
    await this.renderComponent();

    assert.dom(overview.optInBanner.container).doesNotExist();
    assert.dom(cta.button).exists();

    assert
      .dom(overview.createDestination)
      .doesNotExist('create new destination link is hidden if there are no pre-existing destinations');
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

  // WITH DESTINATIONS
  test('it should show the overview cards when destinations exist', async function (assert) {
    this.setup({ isActivated: true, isEnterprise: true, hasFeature: true, hasDestinations: true });
    await this.renderComponent();

    assert.dom(overviewCard.title('Secrets by destination')).exists();
    assert.dom(overviewCard.title('Total destinations')).exists();
    assert.dom(overviewCard.title('Total secrets')).exists();
  });

  test('it should show the table with correct columns, data, badges, and actions', async function (assert) {
    this.setup({ isActivated: true, isEnterprise: true, hasFeature: true, hasDestinations: true });
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
    assert.dom(badge(0)).hasClass('hds-badge--color-success', 'Correct color renders for All synced badge');

    assert.dom(badge(1)).doesNotExist('Status badge does not render for destination with no associations');

    assert.dom(badge(2)).hasText('1 Unsynced', 'Unsynced badge renders');
    assert.dom(badge(2)).hasClass('hds-badge--color-neutral', 'Correct color renders for unsynced badge');

    assert.dom(total(0)).hasText('1', '# of external secrets renders');
    assert.dom(updated(0)).hasText(updatedDate, 'Last updated datetime renders');

    assert.dom(total(1)).hasText('0', '# of external secrets renders for destination with no associations');
    assert
      .dom(updated(1))
      .hasText('—', 'Last updated placeholder renders for destination with no associations');

    await click(actionToggle(0));
    assert.dom(action('sync')).hasText('Sync secrets', 'Sync action renders');
    assert.dom(action('details')).hasText('View synced secrets', 'View synced secrets action renders');
  });

  test('it should paginate the table', async function (assert) {
    this.setup({ isActivated: true, isEnterprise: true, hasFeature: true, hasDestinations: true });
    await this.renderComponent();

    const { name, row } = overview.table;
    assert.dom(row).exists({ count: 3 }, 'Correct number of table rows render based on page size');
    assert.dom(name(0)).hasText('destination-aws', 'First destination renders on page 1');

    await click(PAGE.nextPage);
    await settled();
    assert.dom(overview.table.row).exists({ count: 3 }, 'New items are fetched and rendered on page change');
    assert.dom(name(0)).hasText('destination-gcp', 'First destination renders on page 2');
  });

  test('it should show an empty state if there is an error fetching associations', async function (assert) {
    this.setup({ isActivated: true, isEnterprise: true, hasFeature: true, hasDestinations: true });
    this.server.get('/sys/sync/destinations/:type/:name/associations', () => {
      return new Response(403, {}, { errors: ['Permission denied'] });
    });

    await this.renderComponent();

    assert.dom(emptyStateTitle).hasText('Error fetching information', 'Empty state title renders');
    assert
      .dom(emptyStateMessage)
      .hasText('Ensure that the policy has access to read sync associations.', 'Empty state message renders');
  });

  test('it should show the Totals cards', async function (assert) {
    this.setup({ isActivated: true, isEnterprise: true, hasFeature: true, hasDestinations: true });
    await this.renderComponent();

    const { title, description, actionLink, content } = overviewCard;
    const cardData = [
      {
        cardTitle: 'Total destinations',
        subText: 'The total number of connected destinations.',
        actionText: 'Create new',
        count: '6',
      },
      {
        cardTitle: 'Total secrets',
        subText:
          'The total number of secrets that have been synced from Vault over time. One secret will be counted as one sync client.',
        actionText: 'View billing',
        count: '7',
      },
    ];

    cardData.forEach(({ cardTitle, subText, actionText, count }) => {
      assert.dom(title(cardTitle)).hasText(cardTitle, `${cardTitle} card title renders`);
      assert.dom(description(cardTitle)).hasText(subText, ` ${cardTitle} card description renders`);
      assert.dom(content(cardTitle)).hasText(count, 'Total count renders');
      assert.dom(actionLink(cardTitle)).hasText(actionText, 'Card action renders');
    });
  });
});
