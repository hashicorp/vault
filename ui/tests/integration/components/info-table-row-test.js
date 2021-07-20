import { module, test } from 'qunit';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const VALUE = 'testing';
const LABEL = 'item';
const TYPE = 'array';

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

module('Integration | Component | InfoTableItem', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('value', VALUE);
    this.set('label', LABEL);
    this.set('type', TYPE);
    this.owner.register('service:router', routerService);
    this.router = this.owner.lookup('service:router');
  });

  hooks.afterEach(function() {
    this.owner.unregister('service:store');
  });

  test('it renders', async function(assert) {
    this.set('alwaysRender', true);

    await render(hbs`<InfoTableRow
        @value={{value}}
        @label={{label}}
      />`);

    assert.dom('[data-test-component="info-table-row"]').exists();
    let string = document.querySelector('code').textContent;
    assert.equal(string, VALUE, 'renders value as passed through');

    this.set('value', '');
    assert.dom('[data-test-label-div]').doesNotExist('does not render if no value and alwaysRender is false');
  });

  test('it renders a string with no link if isLink is true and the item type is not an array.', async function(assert) {
    // This could be changed in the component so that it adds a link for any item type, but right it should only add a link if item type is an array.
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

  test('it renders a dash (-) if a label and/or value do not exist', async function(assert) {
    this.set('value', VALUE);
    this.set('label', '');

    await render(hbs`<InfoTableRow
      @value={{value}}
      @label={{label}}
      @alwaysRender={{true}}
    />`);
    assert.dom('[data-test-label-div]').exists('renders label div');
    assert.dom('[data-test-value-div]').exists('renders value div');
    assert.dom('div.column span').hasClass('hs-icon-s', 'Renders a dash (-) for the label');

    this.set('value', '');
    this.set('label', LABEL);
    assert.dom('div.column.is-flex span').hasClass('hs-icon-s', 'Renders a dash (-) for the value');

    this.set('value', '');
    this.set('label', '');
    let dashCount = document.querySelectorAll('.hs-icon-s').length;
    assert.equal(dashCount, 2, 'Renders dash (-) when both label and value do not exist');
  });
});
