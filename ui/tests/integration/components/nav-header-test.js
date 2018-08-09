import { moduleForComponent, test } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import navHeader from 'vault/tests/pages/components/nav-header';

const component = create(navHeader);

moduleForComponent('nav-header', 'Integration | Component | nav header', {
  integration: true,

  beforeEach() {
    component.setContext(this);
  },

  afterEach() {
    component.removeContext();
  },
});

test('it renders', function(assert) {
  this.render(hbs`
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
