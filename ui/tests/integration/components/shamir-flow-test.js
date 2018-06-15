import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import Ember from 'ember';

let response = {
  progress: 1,
  required: 3,
  complete: false,
};

let percent = () => {
  const percent = response.progress / response.required * 100;
  return percent.toFixed(4);
};

let adapter = {
  foo() {
    return Ember.RSVP.resolve(response);
  },
};

const storeStub = Ember.Service.extend({
  adapterFor() {
    return adapter;
  },
});

moduleForComponent('shamir-flow', 'Integration | Component | shamir flow', {
  integration: true,
  beforeEach: function() {
    this.register('service:store', storeStub);
    this.inject.service('store', { as: 'storeService' });
  },
});

test('it renders', function(assert) {
  this.render(hbs`{{shamir-flow formText="like whoa"}}`);

  assert.equal(this.$('form p').text().trim(), 'like whoa', 'renders formText inline');

  this.render(hbs`
    {{#shamir-flow formText="like whoa"}}
      <p>whoa again</p>
    {{/shamir-flow}}
  `);

  assert.equal(this.$('.shamir-progress').length, 0, 'renders no progress bar for no progress');
  assert.equal(this.$('form p').text().trim(), 'whoa again', 'renders the block, not formText');

  this.render(hbs`
    {{shamir-flow progress=1 threshold=5}}
  `);

  assert.ok(this.$('.shamir-progress').text().includes('1/5 keys provided'), 'displays textual progress');

  this.set('errors', ['first error', 'this is fine']);
  this.render(hbs`
    {{shamir-flow errors=errors}}
  `);
  assert.equal(this.$('.message.is-danger').length, 2, 'renders errors');
});

test('it sends data to the passed action', function(assert) {
  this.set('key', 'foo');
  this.render(hbs`
    {{shamir-flow key=key action='foo' thresholdPath='required'}}
  `);
  this.$('[data-test-shamir-submit]').click();
  assert.ok(
    this.$('.shamir-progress').text().includes(`${response.progress}/${response.required} keys provided`),
    'displays the correct progress'
  );
});

test('it checks onComplete to call onShamirSuccess', function(assert) {
  this.set('key', 'foo');
  this.set('onSuccess', function() {
    assert.ok(true, 'onShamirSuccess called');
  });

  this.set('checkComplete', function() {
    assert.ok(true, 'onComplete called');
    // return true so we trigger success call
    return true;
  });

  this.render(hbs`
    {{shamir-flow key=key action='foo' isComplete=(action checkComplete) onShamirSuccess=(action onSuccess)}}
  `);
  this.$('[data-test-shamir-submit]').click();
});

test('it fetches progress on init when fetchOnInit is true', function(assert) {
  this.render(hbs`
    {{shamir-flow action='foo' fetchOnInit=true}}
  `);
  assert.ok(
    this.$('.shamir-progress').text().includes(`${response.progress}/${response.required} keys provided`),
    'displays the correct progress'
  );
});
