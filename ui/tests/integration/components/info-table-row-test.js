import { module, test } from 'qunit';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled, triggerEvent } from '@ember/test-helpers';
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

module('Integration | Component | InfoTableRow', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('value', VALUE);
    this.set('label', LABEL);
    this.set('type', TYPE);
    this.set('default', DEFAULT);
    this.owner.register('service:router', routerService);
    this.router = this.owner.lookup('service:router');
  });

  hooks.afterEach(function() {
    this.owner.unregister('service:store');
  });

  test('it renders', async function(assert) {
    await render(hbs`<InfoTableRow
        @value={{value}}
        @label={{label}}
        @defaultShown={{default}}
      />`);

    assert.dom('[data-test-component="info-table-row"]').exists();
    let string = document.querySelector('code').textContent;
    assert.equal(string, VALUE, 'renders value as passed through');

    this.set('value', '');
    assert
      .dom('[data-test-label-div]')
      .doesNotExist('does not render if no value and alwaysRender is false (even if default exists)');
  });

  test('it renders a tooltip', async function(assert) {
    this.set('tooltipText', 'Tooltip text!');

    await render(hbs`<InfoTableRow
        @value={{value}}
        @label={{label}}
        @tooltipText={{tooltipText}}
      />`);

    await triggerEvent('[data-test-value-div="test label"] .ember-basic-dropdown-trigger', 'mouseenter');
    await settled();
    let tooltip = document.querySelector('div.box').textContent.trim();
    assert.equal(tooltip, 'Tooltip text!', 'renders tooltip text');
  });

  test('it renders a string with no link if isLink is true and the item type is not an array.', async function(assert) {
    // This could be changed in the component so that it adds a link for any item type, but right now it should only add a link if item type is an array.
    await render(hbs`<InfoTableRow
        @value={{value}}
        @label={{label}}
        @isLink={{true}}
      />`);
    let string = document.querySelector('code').textContent;
    assert.equal(string, VALUE, 'renders value in code element and not in a tag');
  });

  test('it renders links if isLink is true and type is array', async function(assert) {
    this.set('valueArray', ['valueArray']);
    await render(hbs`<InfoTableRow
      @value={{valueArray}}
      @label={{label}}
      @isLink={{true}}
      @type={{type}}
    />`);

    assert.dom('[data-test-item="array"]').hasText('valueArray', 'Confirm link with item value exist');
  });

  test('it renders as expected if a label and/or value do not exist', async function(assert) {
    this.set('value', VALUE);
    this.set('label', '');
    this.set('default', '');

    await render(hbs`<InfoTableRow
      @value={{value}}
      @label={{label}}
      @alwaysRender={{true}}
      @defaultShown={{default}}
    />`);

    assert.dom('div.column span').hasClass('hs-icon-s', 'Renders a dash (-) for the label');

    this.set('value', '');
    this.set('label', LABEL);
    assert.dom('div.column.is-flex span').hasClass('hs-icon-s', 'Renders a dash (-) for empty string value');

    this.set('value', null);
    assert.dom('div.column.is-flex span').hasClass('hs-icon-s', 'Renders a dash (-) for null value');

    this.set('value', undefined);
    assert.dom('div.column.is-flex span').hasClass('hs-icon-s', 'Renders a dash (-) for undefined value');

    this.set('default', DEFAULT);
    assert.dom('[data-test-value-div]').hasText(DEFAULT, 'Renders default text if value is empty');

    this.set('value', '');
    this.set('label', '');
    this.set('default', '');
    let dashCount = document.querySelectorAll('.hs-icon-s').length;
    assert.equal(dashCount, 2, 'Renders dash (-) when both label and value do not exist (and no defaults)');
  });

  test('block content overrides any passed in value content', async function(assert) {
    await render(hbs`<InfoTableRow
      @value={{value}}
      @label={{label}}
      @alwaysRender={{true}}>
      Block content is here 
      </InfoTableRow>`);

    let block = document.querySelector('[data-test-value-div]').textContent.trim();
    assert.equal(block, 'Block content is here', 'renders block passed through');
  });
});
