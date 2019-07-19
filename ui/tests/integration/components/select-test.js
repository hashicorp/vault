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
    this.set('name', 'foo');
  });

  test('it renders', async function(assert) {
    await render(hbs`<Select @options={{options}} @label={{label}} @name={{name}}/>`);
    assert.dom('[data-test-select-label]').hasText('Boop');
    assert.dom('[data-test-select="foo"]').exists();
  });

  test('it renders when options is an array of strings', async function(assert) {
    await render(hbs`<Select @options={{options}} @label={{label}} @name={{name}}/>`);

    assert.dom('[data-test-select="foo"]').hasValue('foo');
    assert.equal(this.element.querySelector('[data-test-select="foo"]').options.length, 3);
  });

  test('it renders when options is an array of objects', async function(assert) {
    const objectOptions = [{ value: 'berry', label: 'Berry' }, { value: 'cherry', label: 'Cherry' }];
    this.set('options', objectOptions);
    await render(hbs`<Select @options={{options}} @label={{label}} @name={{name}}/>`);

    assert.dom('[data-test-select="foo"]').hasValue('berry');
    assert.equal(this.element.querySelector('[data-test-select="foo"]').options.length, 2);
  });

  test('it renders when options is an array of custom objects', async function(assert) {
    const objectOptions = [{ day: 'mon', fullDay: 'Monday' }, { day: 'tues', fullDay: 'Tuesday' }];
    const selectedItem = objectOptions[1].day;
    this.setProperties({
      options: objectOptions,
      valueAttribute: 'day',
      labelAttribute: 'fullDay',
      selectedItem: 'tues',
    });

    await render(
      hbs`
        <Select
          @options={{options}}
          @label={{label}}
          @name={{name}}
          @valueAttribute={{valueAttribute}}
          @labelAttribute={{labelAttribute}}
          @selectedItem={{selectedItem}}/>`
    );

    assert.dom('[data-test-select="foo"]').hasValue(selectedItem, 'sets selectedItem by default');
    assert.equal(
      this.element.querySelector('[data-test-select="foo"]').options[1].textContent.trim(),
      'Tuesday',
      'uses the labelAttribute to determine the label'
    );
  });

  test('it renders the selectedItem as selected by default', async function(assert) {
    this.set('selectedItem', 'baz');
    await render(hbs`<Select @options={{options}} @name={{name}} @selectedItem={{selectedItem}}/>`);

    assert.dom('[data-test-select="foo"]').hasValue('baz');
  });

  test('it calls onChange when an item is selected', async function(assert) {
    this.set('onChange', sinon.spy());
    await render(hbs`<Select @options={{options}} @name={{name}} @onChange={{onChange}}/>`);
    await fillIn('[data-test-select="foo"]', 'bar');

    assert.ok(this.onChange.calledOnce);
  });
});
