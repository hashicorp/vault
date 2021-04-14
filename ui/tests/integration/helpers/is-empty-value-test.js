import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const template = hbs`
{{#if (is-empty-value inputValue)}}
Empty
{{else}}
Full
{{/if}}
`;

const emptyObject = {};

const nonEmptyObject = { thing: 0 };

module('Integration | Helper | is-empty-value', function(hooks) {
  setupRenderingTest(hooks);

  test('it is truthy if the value evaluated is undefined', async function(assert) {
    this.set('inputValue', undefined);

    await render(template);

    assert.equal(this.element.textContent.trim(), 'Empty');
  });

  test('it is truthy if the value evaluated is an empty string', async function(assert) {
    this.set('inputValue', '');

    await render(template);

    assert.equal(this.element.textContent.trim(), 'Empty');
  });

  test('it is truthy if the value evaluated is an empty object', async function(assert) {
    this.set('inputValue', emptyObject);

    await render(template);

    assert.equal(this.element.textContent.trim(), 'Empty');
  });
  test('it is falsy if the value evaluated is not an empty object', async function(assert) {
    this.set('inputValue', nonEmptyObject);

    await render(template);

    assert.equal(this.element.textContent.trim(), 'Full');
  });
});
