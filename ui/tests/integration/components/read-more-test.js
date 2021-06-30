import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | read-more', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<ReadMore />`);
    assert.equal(this.element.textContent.trim(), '');

    this.set(
      'description',
      'My super long template block text My super long template block text My super long template block text My super long template block text My super long template block text '
    );
    // Template block usage:
    await render(hbs`
      <div style="width: 50px">
        <ReadMore>
          {{description}}
        </ReadMore>
      </div>
    `);

    assert.equal(this.element.textContent.trim(), description);
    assert.dom('[data-test-readmore-toggle]').exists('toggle exists');
    assert.dom('[data-test-readmore-toggle]').hasText('See more', 'Toggle should have text to see more');
    await click('[data-test-readmore-toggle]');
    assert.dom('[data-test-readmore-toggle]').hasText('See less', 'Toggle should have text to see less');
  });
});
