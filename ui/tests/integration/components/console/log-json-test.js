import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | console/log json', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.codeMirror = this.owner.lookup('service:code-mirror');
  });

  test('it renders', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.on('myAction', function(val) { ... });
    const objectContent = { one: 'two', three: 'four', seven: { five: 'six' }, eight: [5, 6] };
    const expectedText = JSON.stringify(objectContent, null, 2);

    this.set('content', objectContent);

    await render(hbs`{{console/log-json content=content}}`);
    const instance = this.codeMirror.instanceFor(find('[data-test-component=json-editor]').id);

    assert.equal(instance.getValue(), expectedText);
  });
});
