/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { setupRenderingTest } from 'ember-qunit';
import { render, triggerEvent } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const VALUE = 'test value';
const LABEL = 'test label';
const TYPE = 'array';
const DEFAULT = 'some default value';

const routerService = Service.extend({
  transitionTo() {
    return {
      followRedirects() {
        return resolve();
      },
    };
  },
  replaceWith() {
    return resolve();
  },
});

module('Integration | Component | InfoTableRow', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('value', VALUE);
    this.set('label', LABEL);
    this.set('type', TYPE);
    this.set('default', DEFAULT);
    this.owner.register('service:router', routerService);
    this.router = this.owner.lookup('service:router');
  });

  hooks.afterEach(function () {
    this.owner.unregister('service:store');
  });

  test('it renders', async function (assert) {
    await render(hbs`<InfoTableRow
        @value={{this.value}}
        @label={{this.label}}
        @defaultShown={{this.default}}
      />`);

    assert.dom('[data-test-component="info-table-row"]').exists();
    assert.dom('[data-test-row-value]').hasText(VALUE, 'renders value as passed through');

    this.set('value', '');
    assert
      .dom('[data-test-label-div]')
      .doesNotExist('does not render if no value and alwaysRender is false (even if default exists)');
  });

  test('it renders a tooltip', async function (assert) {
    this.set('tooltipText', 'Tooltip text!');

    await render(hbs`<InfoTableRow
        @value={{this.value}}
        @label={{this.label}}
        @tooltipText={{this.tooltipText}}
      />`);

    await triggerEvent('[data-test-value-div="test label"] .ember-basic-dropdown-trigger', 'mouseenter');

    const tooltip = document.querySelector('div.box').textContent.trim();
    assert.strictEqual(tooltip, 'Tooltip text!', 'renders tooltip text');
  });

  test('it should copy tooltip', async function (assert) {
    assert.expect(3);

    this.set('isCopyable', false);

    await render(hbs`
      <InfoTableRow
        @label={{this.label}}
        @value={{this.value}}
        @tooltipText="Foo bar"
        @isTooltipCopyable={{this.isCopyable}}
      />
    `);

    await triggerEvent('[data-test-value-div="test label"] .ember-basic-dropdown-trigger', 'mouseenter');

    assert.dom('[data-test-tooltip-copy]').doesNotExist('Tooltip has no copy button');
    this.set('isCopyable', true);
    assert.dom('[data-test-tooltip-copy]').exists('Tooltip has copy button');
    assert
      .dom('[data-test-tooltip-copy]')
      .hasAttribute('data-test-tooltip-copy', 'Foo bar', 'Copy button will copy the tooltip text');
  });

  test('it renders a string with no link if isLink is true and the item type is not an array.', async function (assert) {
    // This could be changed in the component so that it adds a link for any item type, but right now it should only add a link if item type is an array.
    await render(hbs`<InfoTableRow
        @value={{this.value}}
        @label={{this.label}}
        @isLink={{true}}
      />`);
    assert.dom('[data-test-row-value]').hasText(VALUE, 'renders value in code element and not in a tag');
  });

  test('it renders links if isLink is true and type is array', async function (assert) {
    this.set('valueArray', ['valueArray']);
    await render(hbs`<InfoTableRow
      @value={{this.valueArray}}
      @label={{this.label}}
      @isLink={{true}}
      @type={{this.type}}
    />`);

    assert.dom('[data-test-item="valueArray"]').hasText('valueArray', 'Confirm link with item value exist');
  });

  test('it renders as expected if a label and/or value do not exist', async function (assert) {
    this.set('value', VALUE);
    this.set('label', '');
    this.set('default', '');

    await render(hbs`<InfoTableRow
      @value={{this.value}}
      @label={{this.label}}
      @alwaysRender={{true}}
      @defaultShown={{this.default}}
    />`);

    assert.dom('div.column.is-one-quarter .flight-icon').exists('Renders a dash (-) for the label');

    this.set('value', '');
    this.set('label', LABEL);
    assert.dom('div.column.is-flex-center .flight-icon').exists('Renders a dash (-) for empty string value');

    this.set('value', null);
    assert.dom('div.column.is-flex-center .flight-icon').exists('Renders a dash (-) for null value');

    this.set('value', undefined);
    assert.dom('div.column.is-flex-center .flight-icon').exists('Renders a dash (-) for undefined value');

    this.set('default', DEFAULT);
    assert.dom('[data-test-value-div]').hasText(DEFAULT, 'Renders default text if value is empty');

    this.set('value', '');
    this.set('label', '');
    this.set('default', '');
    const dashCount = document.querySelectorAll('.flight-icon').length;
    assert.strictEqual(
      dashCount,
      2,
      'Renders dash (-) when both label and value do not exist (and no defaults)'
    );
  });

  test('block content overrides any passed in value content', async function (assert) {
    await render(hbs`<InfoTableRow
      @value={{this.value}}
      @label={{this.label}}
      @alwaysRender={{true}}>
      Block content is here
      </InfoTableRow>`);

    const block = document.querySelector('[data-test-value-div]').textContent.trim();
    assert.strictEqual(block, 'Block content is here', 'renders block passed through');
  });

  test('Row renders when block content even if alwaysRender = false', async function (assert) {
    await render(hbs`<InfoTableRow
      @label={{this.label}}
      @alwaysRender={{false}}>
      Block content
    </InfoTableRow>`);
    assert.dom('[data-test-value-div]').exists('renders block');
    assert.dom('[data-test-value-div]').hasText('Block content', 'renders block');
  });

  test('Row does not render empty block content when alwaysRender = false', async function (assert) {
    await render(hbs`<InfoTableRow
      @label={{this.label}}
      @alwaysRender={{false}} />`);
    assert.dom('[data-test-component="info-table-row"]').doesNotExist();
  });

  test('Has dashed label if none provided', async function (assert) {
    await render(hbs`<InfoTableRow
        @value={{this.value}}
      />`);
    assert.dom('[data-test-component="info-table-row"]').exists();
    assert.dom('[data-test-icon="minus"]').exists('renders dash when no label');
  });
  test('Truncates the label if too long', async function (assert) {
    this.set('label', 'abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz');
    await render(hbs`<div style="width: 100px;">
      <InfoTableRow
        @label={{this.label}}
        @value={{this.value}}
      />
    </div>`);
    assert.dom('[data-test-component="info-table-row"]').exists('Row renders');
    assert.dom('[data-test-label-div].label-overflow').exists('Label has class label-overflow');
    await triggerEvent('[data-test-row-label]', 'mouseenter');
    assert.dom('[data-test-label-tooltip]').exists('Label tooltip exists on hover');
  });
  test('Renders if block value and alwaysrender=false', async function (assert) {
    await render(hbs`<InfoTableRow @alwaysRender={{false}}>{{this.value}}</InfoTableRow>`);
    assert.dom('[data-test-component="info-table-row"]').exists();
  });
  test('Does not render if value is empty and alwaysrender=false', async function (assert) {
    await render(hbs`<InfoTableRow @alwaysRender={{false}} @value="" />`);
    assert.dom('[data-test-component="info-table-row"]').doesNotExist();
  });
  test('Renders dash for value if value empty and alwaysRender=true', async function (assert) {
    await render(hbs`<InfoTableRow
        @label={{this.label}}
        @alwaysRender={{true}}
      />`);
    assert.dom('[data-test-component="info-table-row"]').exists();
    assert.dom('[data-test-value-div] [data-test-icon="minus"]').exists('renders dash for value');
  });
  test('Renders block over @value or @defaultShown', async function (assert) {
    await render(hbs`<InfoTableRow
        @label={{this.label}}
        @value="bar"
        @defaultShown="baz"
      >
        foo
      </InfoTableRow>`);
    assert.dom('[data-test-component="info-table-row"]').exists();
    assert.dom('[data-test-value-div]').hasText('foo', 'renders block value');
  });
  test('Renders icons if value is boolean', async function (assert) {
    this.set('value', true);
    await render(hbs`<InfoTableRow
      @label={{this.label}}
      @value={{this.value}}
    />`);

    assert.dom('[data-test-boolean-true]').exists('check icon exists');
    assert.dom('[data-test-value-div]').hasText('Yes', 'Renders yes text');
    this.set('value', false);
    assert.dom('[data-test-boolean-false]').exists('x icon exists');
    assert.dom('[data-test-value-div]').hasText('No', 'renders no text');
  });
  test('Renders data-test attrs passed from parent', async function (assert) {
    this.set('value', true);
    await render(hbs`<InfoTableRow
      @label={{this.label}}
      @value={{this.value}}
      data-test-foo-bar
    />`);

    assert.dom('[data-test-foo-bar]').exists();
  });

  test('Formats the value as date when formatDate present', async function (assert) {
    this.set('value', new Date('2018-04-03T14:15:30'));
    await render(hbs`<InfoTableRow
      @label={{this.label}}
      @value={{this.value}}
      @formatDate={{'yyyy'}}
    />`);

    assert.dom('[data-test-value-div]').hasText('2018', 'Renders date with passed format');
  });

  test('Formats the value as TTL when formatTtl present', async function (assert) {
    this.set('value', 6000);
    await render(hbs`<InfoTableRow
      @label={{this.label}}
      @value={{this.value}}
      @formatTtl={{true}}
    />`);

    assert
      .dom('[data-test-value-div]')
      .hasText('1 hour 40 minutes', 'Translates number value to largest unit with carryover of minutes');
  });

  test('Formats string value when formatTtl present', async function (assert) {
    this.set('value', '45m');
    await render(hbs`<InfoTableRow
      @label={{this.label}}
      @value={{this.value}}
      @formatTtl={{true}}
    />`);

    assert.dom('[data-test-value-div]').hasText('45 minutes', 'it formats string duration');
  });
});
