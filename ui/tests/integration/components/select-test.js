import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

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

  test('it calls onChange when an item is selected', async function(assert) {
    this.set('onChange', sinon.spy());
    await render(hbs`<Select @options={{options}} @onChange={{onChange}}/>`);
    await fillIn('[data-test-select]', 'bar');

    assert.ok(this.onChange.calledOnce);
  });
});
