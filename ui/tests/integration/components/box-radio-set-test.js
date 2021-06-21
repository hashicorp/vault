import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | box-radio-set', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders without a title', async function(assert) {
    await render(hbs`<BoxRadioSet />`);

    assert.equal(this.element.textContent.trim(), '');

    await render(hbs`
      <BoxRadioSet>
        template block text
      </BoxRadioSet>
    `);

    assert.equal(this.element.textContent.trim(), 'template block text');
  });

  test('it renders with a title', async function(assert) {
    this.set('title', 'This Category');

    await render(hbs`<BoxRadioSet @title={{title}} />`);
    assert.equal(this.element.textContent.trim(), 'This Category');

    await render(hbs`
      <BoxRadioSet @title={{title}}>
        <div data-test-inner>template block text</div>
      </BoxRadioSet>
    `);
    assert.dom('[data-test-inner]').hasText('template block text');
    assert.dom('[data-test-box-radio-set-title]').hasText('This Category');
  });
});
