/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

// Minimal row used across most tests — extend per test as needed.
const ROW = { id: 'row-1', name: 'my-entity', status: 'active', path: 'secret/', icon: 'user' };

module('Integration | Component | list-table/cell', function (hooks) {
  setupRenderingTest(hooks);

  // ── valueType: "link" ─────────────────────────────────────────────────────

  module('valueType: "link"', function () {
    test('renders Hds::Link::Inline when no icon is configured', async function (assert) {
      // Validates the most common link cell — a plain text hyperlink navigating
      // to a static route using the row id as the model.
      this.column = { key: 'name', label: 'Name', valueType: 'link', route: 'vault.cluster' };
      this.row = ROW;
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(GENERAL.linkTo(ROW.name)).exists('renders a link with the cell value as text');
      assert.dom(GENERAL.linkTo(ROW.name)).hasTagName('a', 'Hds::Link::Inline renders an anchor element');
    });

    test('renders Hds::Link::Standalone when icon is set via column.icon', async function (assert) {
      // Validates the icon-link path — used for primary name columns that also
      // display an engine or resource type icon alongside the link text.
      this.column = {
        key: 'name',
        label: 'Name',
        valueType: 'link',
        route: 'vault.cluster',
        icon: 'server',
      };
      this.row = ROW;
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(GENERAL.linkTo(ROW.name)).exists('renders a link when a static icon is configured');
    });

    test('renders Hds::Link::Standalone when icon is set via column.iconKey on the row', async function (assert) {
      // Validates the dynamic icon path — the icon name comes from a row field,
      // as used for secrets engine type columns where each engine has a different icon.
      this.column = {
        key: 'name',
        label: 'Name',
        valueType: 'link',
        route: 'vault.cluster',
        iconKey: 'icon',
      };
      this.row = ROW; // ROW.icon = 'user'
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert
        .dom(GENERAL.linkTo(ROW.name))
        .exists('renders a link when icon name is sourced from row[column.iconKey]');
    });

    test('renders nothing when value is falsy', async function (assert) {
      // Guards against rendering a broken link element for rows where the
      // display value is null or an empty string.
      this.column = { key: 'name', label: 'Name', valueType: 'link', route: 'vault.cluster' };
      this.row = ROW;
      this.value = null;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom('a').doesNotExist('no anchor element rendered when value is falsy');
    });

    test('resolves route from routeKey on the row when provided', async function (assert) {
      // Validates the dynamic-route path used when different rows navigate to
      // different routes (e.g. secrets engine list where each engine type has
      // its own overview route).
      this.column = {
        key: 'name',
        label: 'Name',
        valueType: 'link',
        routeKey: 'backendLink',
      };
      this.row = { ...ROW, backendLink: 'vault.cluster.secrets.backend.list' };
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert
        .dom(GENERAL.linkTo(ROW.name))
        .exists('renders a link when route is sourced from row[column.routeKey]');
    });

    test('renders a tooltip when no route is available', async function (assert) {
      this.column = {
        key: 'name',
        label: 'Name',
        valueType: 'link',
        routeKey: 'backendLink',
        unavailableTooltip: 'The UI only supports configuration views for these secret engines.',
      };
      this.row = ROW;
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert
        .dom('[aria-label="The UI only supports configuration views for these secret engines."]')
        .exists();
      assert.dom(GENERAL.linkTo(ROW.name)).doesNotExist('does not render a link without a route');
    });
  });

  // ── valueType: "icon-text" ────────────────────────────────────────────────

  module('valueType: "icon-text"', function () {
    test('renders an icon and text side by side', async function (assert) {
      // Validates the basic icon-text cell used for engine type columns — an
      // icon followed by a display label in the same flex row.
      this.column = { key: 'name', label: 'Type', valueType: 'icon-text', icon: 'database' };
      this.row = ROW;
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(GENERAL.icon()).exists('renders an icon element');
      assert.dom('span').hasText(ROW.name, 'renders the cell value as text');
    });

    test('reads icon name from row[column.iconKey] when set', async function (assert) {
      // Validates that per-row icon names override the static column.icon fallback.
      this.column = { key: 'name', label: 'Type', valueType: 'icon-text', iconKey: 'icon' };
      this.row = ROW; // ROW.icon = 'user'
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(GENERAL.icon('user')).exists('uses the icon name from row[iconKey]');
    });

    test('reads display text from row[column.textKey] when set', async function (assert) {
      // Validates the textKey override — used when the display label is a
      // different field than the column key (e.g. key:"type", textKey:"displayName").
      this.column = {
        key: 'name',
        label: 'Type',
        valueType: 'icon-text',
        icon: 'database',
        textKey: 'path',
      };
      this.row = ROW; // ROW.path = 'secret/'
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom('span').hasText(ROW.path, 'renders row[textKey] as the display text instead of value');
    });

    test('wraps icon in a tooltip button when column.tooltip is set', async function (assert) {
      // Validates tooltip rendering — used for accessor columns and other cells
      // where the icon meaning needs disambiguation.
      this.column = {
        key: 'name',
        label: 'Type',
        valueType: 'icon-text',
        icon: 'info',
        tooltip: 'A unique identifier for this engine',
      };
      this.row = ROW;
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert
        .dom('button[aria-label="A unique identifier for this engine"]')
        .exists('tooltip button rendered with aria-label from column.tooltip');
    });

    test('falls back to "lock" icon when no icon config is provided', async function (assert) {
      // Validates the fallback icon — the template hardcodes "lock" when neither
      // iconKey nor icon is set, so the cell always renders a visible icon.
      this.column = { key: 'name', label: 'Type', valueType: 'icon-text' };
      this.row = { id: 'row-1', name: 'my-entity' }; // no icon field
      this.value = 'my-entity';

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(GENERAL.icon('lock')).exists('falls back to the "lock" icon when none is configured');
    });
  });

  // ── valueType: "copy" ─────────────────────────────────────────────────────

  module('valueType: "copy"', function () {
    test('renders a copy button and the plain value text', async function (assert) {
      // Validates that accessor-style copy cells render both the copy trigger
      // and the raw value for visual reference.
      this.column = { key: 'id', label: 'Accessor', valueType: 'copy' };
      this.row = ROW;
      this.value = ROW.id;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(GENERAL.copyButton).exists('renders the copy button');
      assert.dom(this.element).includesText(ROW.id, 'renders the plain value text alongside the copy button');
    });

    test('renders nothing when value is falsy', async function (assert) {
      // Guards against rendering an empty copy button for rows where the
      // accessor or token is null.
      this.column = { key: 'id', label: 'Accessor', valueType: 'copy' };
      this.row = ROW;
      this.value = null;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(GENERAL.copyButton).doesNotExist('copy button not rendered when value is falsy');
    });
  });

  // ── valueType: "status" ───────────────────────────────────────────────────

  module('valueType: "status"', function () {
    test('renders an Hds::Badge with color from statusMap', async function (assert) {
      // Validates that the status cell maps the raw value to the configured
      // badge color via statusMap — used for enabled/disabled or active/inactive columns.
      this.column = {
        key: 'status',
        label: 'Status',
        valueType: 'status',
        statusMap: { active: 'success', inactive: 'warning' },
      };
      this.row = ROW; // ROW.status = 'active'
      this.value = 'active';

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(GENERAL.badge('active')).exists('renders a badge with the cell value as text');
      assert
        .dom(GENERAL.badge('active'))
        .hasClass('hds-badge--color-success', 'badge color matches statusMap entry for "active"');
    });

    test('renders a neutral badge for values not in statusMap', async function (assert) {
      // Validates the fallback behavior — unknown status values render a badge
      // without crashing, using the default (neutral) badge color.
      this.column = {
        key: 'status',
        label: 'Status',
        valueType: 'status',
        statusMap: { active: 'success' },
      };
      this.row = { ...ROW, status: 'unknown-status' };
      this.value = 'unknown-status';

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(GENERAL.badge('unknown-status')).exists('renders a badge for unrecognized status values');
    });

    test('renders nothing when value is falsy', async function (assert) {
      this.column = { key: 'status', label: 'Status', valueType: 'status', statusMap: {} };
      this.row = ROW;
      this.value = null;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(GENERAL.badge()).doesNotExist('no badge rendered when value is falsy');
    });
  });

  // ── valueType: "date-time" ────────────────────────────────────────────────

  module('valueType: "date-time"', function () {
    test('renders a <time> element with the formatted date', async function (assert) {
      // Validates that timestamp values are wrapped in a <time> element with
      // the datetime attribute set and the value formatted for display.
      this.column = { key: 'created_at', label: 'Created', valueType: 'date-time' };
      this.row = { id: 'row-1', created_at: '2024-06-01T12:00:00.000Z' };
      this.value = '2024-06-01T12:00:00.000Z';

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom('time').exists('renders a <time> element');
      assert
        .dom('time')
        .hasAttribute('datetime', '2024-06-01T12:00:00.000Z', 'datetime attribute holds the raw ISO value');
      assert.dom('time').hasText(/Jun 01, 2024/, 'formatted date text is rendered inside the <time> element');
    });

    test('renders nothing when value is falsy', async function (assert) {
      this.column = { key: 'created_at', label: 'Created', valueType: 'date-time' };
      this.row = ROW;
      this.value = null;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom('time').doesNotExist('no <time> element rendered when value is falsy');
    });
  });

  // ── customTableItem yield ─────────────────────────────────────────────────

  module('customTableItem yield', function () {
    test('yields the customTableItem block when column.customTableItem is true', async function (assert) {
      // Validates that the escape hatch for fully custom cell rendering passes
      // the row, column, and value through to the caller's block.
      this.column = { key: 'name', label: 'Name', customTableItem: true };
      this.row = ROW;
      this.value = ROW.name;

      await render(hbs`
        <ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}}>
          <:customTableItem as |row col val|>
            <span data-test-custom-cell>{{val}} — {{col.key}} — {{row.id}}</span>
          </:customTableItem>
        </ListTable::Cell>
      `);

      assert
        .dom('[data-test-custom-cell]')
        .hasText(
          `${ROW.name} — name — ${ROW.id}`,
          'yields row, column, and value to the customTableItem block'
        );
    });
  });

  // ── plain text fallback ───────────────────────────────────────────────────

  module('plain text fallback', function () {
    test('renders the value as plain text when no valueType is set', async function (assert) {
      // Validates the default rendering path used by simple string columns
      // that need no special display treatment.
      this.column = { key: 'name', label: 'Name' };
      this.row = ROW;
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(this.element).includesText(ROW.name, 'renders the raw string value');
    });

    test('JSON-stringifies object values in the plain text fallback', async function (assert) {
      // Validates that nested objects (e.g. metadata maps from an API response)
      // do not cause a silent render failure — they are stringified instead.
      this.column = { key: 'meta', label: 'Meta' };
      this.row = { id: 'row-1', meta: { version: 2 } };
      this.value = { version: 2 };

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(this.element).includesText('{ "version": 2 }', 'object values are JSON-stringified');
    });

    test('renders nothing when the value is falsy', async function (assert) {
      // Null/undefined/empty-string values should produce an empty cell.
      this.column = { key: 'name', label: 'Name' };
      this.row = { name: null };
      this.value = null;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert.dom(this.element).hasText('', 'renders nothing for a falsy value');
    });

    test('renders a TooltipButton when the cell value overflows the parent element', async function (assert) {
      this.column = { key: 'name', label: 'Name' };
      this.row = { name: 'a'.repeat(300) };
      this.value = this.row.name;

      await render(hbs`
        <div style="width:1px; overflow:hidden; white-space:nowrap; display:block;">
          <ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />
        </div>
      `);

      assert
        .dom(GENERAL.tooltip(this.value))
        .exists('tooltip button is shown when the cell value overflows the container');
    });

    test('does not render a TooltipButton when the cell value fits', async function (assert) {
      this.column = { key: 'name', label: 'Name' };
      this.row = ROW;
      this.value = ROW.name;

      await render(hbs`<ListTable::Cell @column={{this.column}} @row={{this.row}} @value={{this.value}} />`);

      assert
        .dom(GENERAL.tooltip(this.value))
        .doesNotExist('no tooltip button when the full value fits in the cell');
    });
  });
});
