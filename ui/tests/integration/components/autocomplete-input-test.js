import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, typeIn, triggerEvent, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | autocomplete-input', function (hooks) {
  setupRenderingTest(hooks);

  test('it should render label', async function (assert) {
    await render(
      hbs`
      <AutocompleteInput
        @label={{this.label}}
        @subText={{this.subText}}
      />`
    );

    assert.dom('label').doesNotExist('Label is hidden when not provided');
    this.setProperties({
      label: 'Some label',
      subText: 'Some description',
    });
    assert.dom('label').hasText('Some label', 'Label renders');
    assert.dom('[data-test-label-subtext]').hasText('Some description', 'Sub text renders');
  });

  test('it should function as standard input', async function (assert) {
    assert.expect(3);
    const changeValue = 'foo bar';
    this.value = 'test';
    this.placeholder = 'text goes here';
    this.onChange = (value) => assert.strictEqual(value, changeValue, 'Value sent in onChange callback');

    await render(
      hbs`
      <AutocompleteInput
        @value={{this.value}}
        @placeholder={{this.placeholder}}
        @onChange={{this.onChange}}
      />`
    );

    assert.dom('input').hasAttribute('placeholder', this.placeholder, 'Input placeholder renders');
    assert.dom('input').hasValue(this.value, 'Initial input value renders');
    await typeIn('input', changeValue);
  });

  test('it should trigger dropdown', async function (assert) {
    await render(
      hbs`
      <AutocompleteInput
        @value={{this.value}}
        @optionsTrigger="$"
        @options={{this.options}}
        @onChange={{fn (mut this.value)}}
      />`
    );

    await typeIn('input', '$');
    await triggerEvent('input', 'input', { data: '$' }); // simulate InputEvent for data prop with character pressed
    assert.dom('.autocomplete-input-option').doesNotExist('Trigger does not open dropdown with no options');

    this.set('options', [
      { label: 'Foo', value: '$foo' },
      { label: 'Bar', value: 'bar' },
    ]);
    await triggerEvent('input', 'input', { data: '$' });
    const options = this.element.querySelectorAll('.autocomplete-input-option');
    options.forEach((o, index) => {
      assert.dom(o).hasText(this.options[index].label, 'Label renders for option');
    });

    await click(options[0]);
    assert.dom('input').isFocused('Focus is returned to input after selecting option');
    assert
      .dom('input')
      .hasValue('$foo', 'Value is updated correctly. Trigger character is not prepended to value.');

    await typeIn('input', '-$');
    await triggerEvent('input', 'input', { data: '$' });
    await click('.autocomplete-input-option:last-child');
    assert
      .dom('input')
      .hasValue('$foo-$bar', 'Value is updated correctly. Trigger character is prepended to option.');
    assert.strictEqual(this.value, '$foo-$bar', 'Value prop is updated correctly onChange');
  });
});
