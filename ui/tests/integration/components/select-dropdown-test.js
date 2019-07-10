import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const OPTIONS = ['foo', 'bar', 'baz'];
const LABEL = 'Boop';

module('Integration | Component | select-dropdown', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('options', OPTIONS);
    this.set('dropdownLabel', LABEL);
  });

  test('it renders with options', async function(assert) {
    await render(hbs`<SelectDropdown @options={{options}} @dropdownLabel={{dropdownLabel}}/>`);

    assert.dom('[data-test-select-label]').hasText('Boop');
    assert.dom('[data-test-select-dropdown]').hasValue('foo', 'shows the first item by default');

    assert.equal(
      this.element.querySelector('[data-test-select-dropdown]').options.length,
      3,
      'it adds an option for each year in the data set'
    );
  });

  test('it renders the selectedItem as selected by default', async function(assert) {
    this.set('selectedItem', 'baz');
    await render(hbs`<SelectDropdown @options={{options}} @selectedItem={{selectedItem}}/>`);

    assert.dom('[data-test-select-dropdown]').hasValue('baz');
  });
});
