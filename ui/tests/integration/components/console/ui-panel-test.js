import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('console/ui-panel', 'Integration | Component | console/ui panel', {
  integration: true
});

test('it renders', function(assert) {

  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });

  this.render(hbs`{{console/ui-panel}}`);

  assert.equal(this.$().text().trim(), '');

  // Template block usage:
  this.render(hbs`
    {{#console/ui-panel}}
      template block text
    {{/console/ui-panel}}
  `);

  assert.equal(this.$().text().trim(), 'template block text');
});
