import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const OPTIONS = ['foo', 'bar', 'baz'];
const LABEL = 'Boop';

module('Integration | Component | Select', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('options', OPTIONS);
    this.set('label', LABEL);
  });

  test('it renders with options', async function(assert) {
    await render(hbs`<Select @options={{options}} @label={{label}}/>`);

    assert.dom('[data-test-select-label]').hasText('Boop');
    assert.dom('[data-test-select]').hasValue('foo', 'shows the first item by default');

    assert.equal(
      this.element.querySelector('[data-test-select]').options.length,
      3,
      'it adds an option for each year in the data set'
    );
  });

  test('it renders the selectedItem as selected by default', async function(assert) {
    this.set('selectedItem', 'baz');
    await render(hbs`<Select @options={{options}} @selectedItem={{selectedItem}}/>`);

    assert.dom('[data-test-select]').hasValue('baz');
  });
});
