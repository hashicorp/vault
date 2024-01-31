/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

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

const {
  title,
  breadcrumb,
  tab,
  overviewCard,
  cta,
  overview,
  pagination,
  emptyStateTitle,
  emptyStateMessage,
} = PAGE;

module('Integration | Component | sync | Page::Overview', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';
    syncScenario(this.server);
    syncHandlers(this.server);

    const store = this.owner.lookup('service:store');
    this.destinations = await store.query('sync/destination', {});

    await render(
      hbs`<Secrets::Page::Overview @destinations={{this.destinations}} @totalAssociations={{7}} />`,
      {
        owner: this.engine,
      }
    );
  });

  test('it should render landing cta component for community', async function (assert) {
    this.version.type = 'community';
    this.set('destinations', []);
    await settled();
    assert.dom(title).hasText('Secrets Sync Enterprise feature', 'Page title renders');
    assert.dom(cta.button).doesNotExist('Create first destination button does not render');
  });

  test('it should render landing cta component for enterprise', async function (assert) {
    this.set('destinations', []);
    await settled();
    assert.dom(title).hasText('Secrets Sync Beta', 'Page title renders');
    assert.dom(cta.button).hasText('Create first destination', 'CTA action renders');
    assert.dom(cta.summary).exists('CTA renders');
  });

  test('it should render header, tabs and toolbar for overview state', async function (assert) {
    assert.dom(title).hasText('Secrets Sync Beta', 'Page title renders');
    assert.dom(breadcrumb).exists({ count: 1 }, 'Correct number of breadcrumbs render');
    assert.dom(breadcrumb).includesText('Secrets Sync', 'Top level breadcrumb renders');
    assert.dom(cta.button).doesNotExist('CTA does not render');
    assert.dom(tab('Overview')).hasText('Overview', 'Overview tab renders');
    assert.dom(tab('Destinations')).hasText('Destinations', 'Destinations tab renders');
    assert.dom(overview.createDestination).hasText('Create new destination', 'Toolbar action renders');
  });

  test('it should render secrets by destination table', async function (assert) {
    const { icon, name, badge, total, updated, actionToggle, action } = overview.table;
    const updatedDate = dateFormat(
      [new Date('2023-09-20T10:51:53.961861096-04:00'), 'MMMM do yyyy, h:mm:ss a'],
      {}
    );
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

    assert.dom(total(0)).hasText('1', '# of secrets renders');
    assert.dom(updated(0)).hasText(updatedDate, 'Last updated datetime renders');

    assert.dom(total(1)).hasText('0', '# of secrets render for destination with no associations');
    assert
      .dom(updated(1))
      .hasText('—', 'Last updated placeholder renders for destination with no associations');

    await click(actionToggle(0));
    assert.dom(action('sync')).hasText('Sync secrets', 'Sync action renders');
    assert.dom(action('details')).hasText('View synced secrets', 'View synced secrets action renders');
  });

  test('it should paginate secrets by destination table', async function (assert) {
    const { name, row } = overview.table;
    assert.dom(row).exists({ count: 3 }, 'Correct number of table rows render based on page size');
    assert.dom(name(0)).hasText('destination-aws', 'First destination renders on page 1');

    await click(pagination.next);
    assert.dom(overview.table.row).exists({ count: 3 }, 'New items are fetched and rendered on page change');
    assert.dom(name(0)).hasText('destination-gcp', 'First destination renders on page 2');
  });

  test('it should display empty state for secrets by destination table', async function (assert) {
    this.server.get('/sys/sync/destinations/:type/:name/associations', () => {
      return new Response(403, {}, { errors: ['Permission denied'] });
    });
    // since the request resolved trigger a page change and return an error from the associations endpoint
    await click(pagination.next);
    assert.dom(emptyStateTitle).hasText('Error fetching information', 'Empty state title renders');
    assert
      .dom(emptyStateMessage)
      .hasText('Ensure that the policy has access to read sync associations.', 'Empty state message renders');
  });

  test('it should render totals cards', async function (assert) {
    const { title, description, action, content } = overviewCard;
    const cardData = [
      {
        cardTitle: 'Total destinations',
        subText: 'The total number of connected destinations',
        actionText: 'Create new',
        count: '6',
      },
      {
        cardTitle: 'Total sync associations',
        subText: 'Total sync associations that count towards client count',
        actionText: 'View billing',
        count: '7',
      },
    ];

    cardData.forEach(({ cardTitle, subText, actionText, count }) => {
      assert.dom(title(cardTitle)).hasText(cardTitle, 'Overview card title renders');
      assert.dom(description(cardTitle)).hasText(subText, 'Destinations overview card description renders');
      assert.dom(action(cardTitle)).hasText(actionText, 'Card action renders');
      assert.dom(content(cardTitle)).hasText(count, 'Total count renders');
    });
  });
});
