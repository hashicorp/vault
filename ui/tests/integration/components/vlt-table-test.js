import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const DATA = [
  {
    foo: 'panda',
    bar: 50,
  },
  {
    foo: 'moose',
    bar: 45,
  },
  {
    foo: 'ocelot',
    bar: 55,
  },
];

module('Integration | Component | vlt-table', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('data', DATA);
  });

  test('it renders with headers and values', async function(assert) {
    await render(hbs`<VltTable @data={{DATA}}/>`);

    assert.ok(this.element.textContent.includes('foo'));
    assert.ok(this.element.textContent.includes('bar'));
    assert.ok(this.element.textContent.includes('moose'));
    assert.ok(this.element.textContent.includes(55));
  });
});
