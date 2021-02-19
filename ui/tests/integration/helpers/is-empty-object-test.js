import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const emptyObject = {};

const nonEmptyObject = { thing: 0 };

module('Integration | Helper | is-empty-object', function(hooks) {
  setupRenderingTest(hooks);

  test('it is truthy if the object evaluated is an empty object', async function(assert) {
    this.set('inputValue', emptyObject);

    await render(hbs`
      {{#if (is-empty-object inputValue)}}
      Empty
      {{else}}
      Full
      {{/if}}
    `);

    assert.equal(this.element.textContent.trim(), 'Empty');
  });
  test('it is falsy if the object evaluated is not an empty object', async function(assert) {
    this.set('inputValue', nonEmptyObject);

    await render(hbs`
      {{#if (is-empty-object inputValue)}}
      Empty
      {{else}}
      Full
      {{/if}}
    `);

    assert.equal(this.element.textContent.trim(), 'Full');
  });
});
