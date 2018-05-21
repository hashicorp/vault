import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('console/log-json', 'Integration | Component | console/log json', {
  integration: true,

  beforeEach() {
    this.inject.service('code-mirror', { as: 'codeMirror' });
  },
});

test('it renders', function(assert) {
  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });
  const objectContent = { one: 'two', three: 'four', seven: { five: 'six' }, eight: [5, 6] };
  const expectedText = JSON.stringify(objectContent, null, 2);

  this.set('content', objectContent);

  this.render(hbs`{{console/log-json content=content}}`);
  const instance = this.codeMirror.instanceFor(this.$('[data-test-component=json-editor]').attr('id'));

  assert.equal(instance.getValue(), expectedText);
});
