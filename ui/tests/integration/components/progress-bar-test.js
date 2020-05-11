import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, pauseTest } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | progress-bar', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('threshold', 4);
    this.set('progress', 3);
  });

  test('it renders', async function(assert) {
    await render(hbs`<ProgressBar @threshold={{threshold}} @progress={{progress}}/>`);

    assert.dom('.progress-bar').exists();
  });
});
