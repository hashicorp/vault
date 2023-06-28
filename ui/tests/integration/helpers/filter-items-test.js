import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Helper | filter-items', function (hooks) {
  setupRenderingTest(hooks);

  test('it returns empty array if falsy value', async function (assert) {
    this.set('inputValue', null);

    await render(hbs`
    {{#each (filter-items this.inputValue "something") as |thing|}}
      <div data-test-thing>hello</div>
    {{/each}}
    `);
    assert.dom('[data-test-thing]').doesNotExist('no items rendered');
  });

  test('it searches on ID by default', async function (assert) {
    this.set('filter', 'titan');
    this.set('inputValue', [
      { id: 'titanic', name: 'something' },
      { id: 'thing', name: 'titanic' },
      { id: 'titan', name: 'something' },
    ]);

    await render(hbs`
    {{#each (filter-items this.inputValue this.filter) as |thing|}}
      <div data-test-thing={{thing.id}}>{{thing.name}}</div>
    {{/each}}
    `);

    assert.dom('[data-test-thing]').exists({ count: 2 }, 'Filters to 2 matching items');
    assert.dom('[data-test-thing="titanic"]').exists();
    assert.dom('[data-test-thing="titan"]').exists();

    this.set('filter', 'titani');
    assert.dom('[data-test-thing]').exists({ count: 1 }, 'Filters to 1 matching item');
    assert.dom('[data-test-thing="titanic"]').exists();
  });

  test('it searches on a passed attribute', async function (assert) {
    this.set('filter', 'foo');
    this.set('inputValue', [
      { id: 'thing-1', name: 'foo' },
      { id: 'foo', name: 'bar' },
      { id: 'thing-3', name: 'foobar' },
    ]);

    await render(hbs`
    {{#each (filter-items this.inputValue this.filter (hash attr="name")) as |thing|}}
      <div data-test-thing={{thing.id}}>{{thing.name}}</div>
    {{/each}}
    `);

    assert.dom('[data-test-thing]').exists({ count: 2 }, 'Filters to 2 matching items');
    assert.dom('[data-test-thing="thing-1"]').hasText('foo');
    assert.dom('[data-test-thing="thing-3"]').hasText('foobar');

    this.set('filter', 'foob');
    assert.dom('[data-test-thing]').exists({ count: 1 }, 'Filters to 1 matching item');
    assert.dom('[data-test-thing="thing-3"]').exists();

    this.set('filter', 'nothing');
    assert.dom('[data-test-thing]').doesNotExist();
  });
});
