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

  test('it renders', async function(assert) {
    await render(hbs`<Select @options={{options}} @label={{label}}/>`);

    assert.dom('[data-test-select-label]').hasText('Boop');
    assert.dom('[data-test-select]').exists();
  });

  test('it renders when options is an array of strings', async function(assert) {
    await render(hbs`<Select @options={{options}} @label={{label}}/>`);

    assert.dom('[data-test-select]').hasValue('foo');
    assert.equal(this.element.querySelector('[data-test-select]').options.length, 3);
  });

  test('it renders when options is an array of objects', async function(assert) {
    const objectOptions = [{ value: 'berry', label: 'Berry' }, { value: 'cherry', label: 'Cherry' }];
    this.set('options', objectOptions);
    await render(hbs`<Select @options={{options}} @label={{label}}/>`);

    assert.dom('[data-test-select]').hasValue('berry');
    assert.equal(this.element.querySelector('[data-test-select]').options.length, 2);
  });

  test('it renders when options is an array of custom objects', async function(assert) {
    const objectOptions = [{ day: 'mon', fullDay: 'Monday' }, { day: 'tues', fullDay: 'Tuesday' }];
    this.setProperties({
      options: objectOptions,
      valueAttribute: 'day',
      labelAttribute: 'fullDay',
    });

    await render(
      hbs`
        <Select
          @options={{options}}
          @label={{label}}
          @valueAttribute={{valueAttribute}}
          @labelAttribute={{labelAttribute}}/>`
    );

    assert.dom('[data-test-select]').hasValue('mon');
    assert.equal(this.element.querySelector('[data-test-select]').options[1].textContent, 'Tuesday');
  });

  test('it calls onChange when an item is selected', async function(assert) {
    this.set('onChange', sinon.spy());
    await render(hbs`<Select @options={{options}} @onChange={{onChange}}/>`);
    await fillIn('[data-test-select]', 'bar');

    assert.ok(this.onChange.calledOnce);
  });
});
