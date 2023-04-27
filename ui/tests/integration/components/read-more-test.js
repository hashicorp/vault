/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | read-more', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<ReadMore />`);
    assert.dom(this.element).hasText('');
  });
  test('it can toggle open and closed when text is longer than parent', async function (assert) {
    this.set(
      'description',
      'My super long template block text dignissim dictum sem, ut varius ligula lacinia quis.'
    );
    await render(hbs`
      <div style="width: 150px">
        <ReadMore>
          {{this.description}}
        </ReadMore>
      </div>
    `);
    assert.dom('[data-test-readmore-content]').includesText(this.description);
    assert.dom('[data-test-readmore-toggle]').exists('toggle exists');
    assert.dom('[data-test-readmore-toggle]').hasText('See More', 'Toggle should have text to see more');
    assert
      .dom('.overflow-ellipsis.is-closed')
      .exists('Overflow div has is-closed class when more text to show');
    await click('[data-test-readmore-toggle] button');
    assert.dom('.overflow-ellipsis').exists('Div with overflow class still exists');
    assert.dom('.overflow-ellipsis.is-closed').doesNotExist('Div with overflow class no longer is-closed');
    assert.dom('[data-test-readmore-toggle]').hasText('See Less', 'Toggle should have text to see less');
  });

  test('it does not show show more button if parent is wider than content', async function (assert) {
    this.set('description', 'Hello');
    await render(hbs`
      <div style="width: 150px">
        <ReadMore>
          {{this.description}}
        </ReadMore>
      </div>
    `);
    assert.dom('[data-test-readmore-content]').includesText(this.description);
    assert.dom('[data-test-readmore-toggle]').doesNotExist('toggle exists');
    assert.dom('.overflow-ellipsis').exists('Overflow div exists');
  });
});
