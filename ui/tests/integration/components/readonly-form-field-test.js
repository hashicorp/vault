import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const minimumAttr = {
  name: 'my-input',
  type: 'text',
};
const customLabelAttr = {
  name: 'test-input',
  type: 'text',
  options: {
    subText: 'Subtext here',
    label: 'Custom-label',
  },
};

module('Integration | Component | readonly-form-field', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    this.set('attr', minimumAttr);
    await render(hbs`<ReadonlyFormField @attr={{attr}} @value="value" />`);
    assert
      .dom('[data-test-readonly-label]')
      .includesText('My input', 'formats the attr name when no label provided');
    assert.dom(`input`).hasValue('value', 'Uses the value as passed');
    assert.dom(`input`).hasAttribute('readonly');
  });

  test('it renders with options', async function(assert) {
    this.set('attr', customLabelAttr);
    await render(hbs`<ReadonlyFormField @attr={{attr}} @value="another value" />`);
    assert
      .dom('[data-test-readonly-label]')
      .includesText('Custom-label', 'Uses the provided label as passed');
    assert.dom('.sub-text').includesText('Subtext here', 'Renders subtext');
    assert.dom('input').hasValue('another value', 'Uses the value as passed');
    assert.dom('input').hasAttribute('readonly');
  });
});
