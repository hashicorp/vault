import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import navHeader from 'vault/tests/pages/components/nav-header';

const component = create(navHeader);

module('Integration | Component | nav header', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
  });

  hooks.afterEach(function() {
    component.removeContext();
  });

  test('it renders', async function(assert) {
    await render(hbs`
        {{#nav-header as |h|}}
          {{#h.home}}
            Home!
          {{/h.home}}
          {{#h.items}}
            Some Items
          {{/h.items}}
          {{#h.main}}
            Main stuff here
          {{/h.main}}
        {{/nav-header}}
      `);

    assert.ok(component.ele, 'renders the outer element');
    assert.equal(component.homeText.trim(), 'Home!', 'renders home content');
    assert.equal(component.itemsText.trim(), 'Some Items', 'renders items content');
    assert.equal(component.mainText.trim(), 'Main stuff here', 'renders items content');
  });
});
