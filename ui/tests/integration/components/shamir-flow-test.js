import { run } from '@ember/runloop';
import Service from '@ember/service';
import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, findAll, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

let response = {
  progress: 1,
  required: 3,
  complete: false,
};

let adapter = {
  foo() {
    return resolve(response);
  },
};

const storeStub = Service.extend({
  adapterFor() {
    return adapter;
  },
});

module('Integration | Component | shamir flow', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeStub);
      this.storeService = this.owner.lookup('service:store');
    });
  });

  test('it renders', async function(assert) {
    await render(hbs`{{shamir-flow formText="like whoa"}}`);

    assert.equal(
      find('form [data-test-form-text]').textContent.trim(),
      'like whoa',
      'renders formText inline'
    );

    await render(hbs`
      {{#shamir-flow formText="like whoa"}}
        <p>whoa again</p>
      {{/shamir-flow}}
    `);

    assert.equal(findAll('.shamir-progress').length, 0, 'renders no progress bar for no progress');
    assert.equal(
      find('form [data-test-form-text]').textContent.trim(),
      'whoa again',
      'renders the block, not formText'
    );

    await render(hbs`
      {{shamir-flow progress=1 threshold=5}}
    `);

    assert.ok(
      find('.shamir-progress').textContent.includes('1/5 keys provided'),
      'displays textual progress'
    );

    this.set('errors', ['first error', 'this is fine']);
    await render(hbs`
      {{shamir-flow errors=errors}}
    `);
    assert.equal(findAll('.message.is-danger').length, 2, 'renders errors');
  });

  test('it sends data to the passed action', async function(assert) {
    this.set('key', 'foo');
    await render(hbs`
      {{shamir-flow key=key action='foo' thresholdPath='required'}}
    `);
    await click('[data-test-shamir-submit]');
    assert.ok(
      find('.shamir-progress').textContent.includes(
        `${response.progress}/${response.required} keys provided`
      ),
      'displays the correct progress'
    );
  });

  test('it checks onComplete to call onShamirSuccess', async function(assert) {
    this.set('key', 'foo');
    this.set('onSuccess', function() {
      assert.ok(true, 'onShamirSuccess called');
    });

    this.set('checkComplete', function() {
      assert.ok(true, 'onComplete called');
      // return true so we trigger success call
      return true;
    });

    await render(hbs`
      {{shamir-flow key=key action='foo' isComplete=(action checkComplete) onShamirSuccess=(action onSuccess)}}
    `);
    await click('[data-test-shamir-submit]');
  });

  test('it fetches progress on init when fetchOnInit is true', async function(assert) {
    await render(hbs`
      {{shamir-flow action='foo' fetchOnInit=true}}
    `);
    assert.ok(
      find('.shamir-progress').textContent.includes(
        `${response.progress}/${response.required} keys provided`
      ),
      'displays the correct progress'
    );
  });
});
