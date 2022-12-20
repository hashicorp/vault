import { run } from '@ember/runloop';
import Service from '@ember/service';
import { resolve } from 'rsvp';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const response = {
  progress: 1,
  required: 3,
  complete: false,
};

const adapter = {
  foo() {
    return resolve(response);
  },
};

const storeStub = Service.extend({
  adapterFor() {
    return adapter;
  },
});

module('Integration | Component | shamir flow', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.foo = function () {};
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeStub);
      this.storeService = this.owner.lookup('service:store');
    });
  });

  test('it renders', async function (assert) {
    await render(hbs`<ShamirFlow @formText="like whoa" />`);

    assert.dom('form [data-test-form-text]').hasText('like whoa', 'renders formText inline');

    await render(hbs`
      <ShamirFlow @formText="like whoa">
        <p>whoa again</p>
      </ShamirFlow>
    `);

    assert.dom('.shamir-progress').doesNotExist('renders no progress bar for no progress');
    assert.dom('form [data-test-form-text]').hasText('whoa again', 'renders the block, not formText');

    await render(hbs`
      <ShamirFlow @progress={{1}} @threshold={{5}} />
    `);
    assert.dom('.shamir-progress').hasText('1/5 keys provided', 'displays textual progress');

    this.set('errors', ['first error', 'this is fine']);
    await render(hbs`
    <ShamirFlow @errors={{this.errors}} />
    `);
    assert.dom('.message.is-danger').exists({ count: 2 }, 'renders errors');
  });

  test('it sends data to the passed action', async function (assert) {
    this.set('key', 'foo');
    await render(hbs`
      <ShamirFlow @key={{this.key}} @action="foo" @thresholdPath="required" />
    `);
    await click('[data-test-shamir-submit]');
    assert
      .dom('.shamir-progress')
      .hasText(`${response.progress}/${response.required} keys provided`, 'displays the correct progress');
  });

  test('it checks onComplete to call onShamirSuccess', async function (assert) {
    assert.expect(2);
    this.set('key', 'foo');
    this.set('onSuccess', function () {
      assert.ok(true, 'onShamirSuccess called');
    });

    this.set('checkComplete', function () {
      assert.ok(true, 'onComplete called');
      // return true so we trigger success call
      return true;
    });

    await render(hbs`
      <ShamirFlow @key={{this.key}} @action="foo" @isComplete={{action this.checkComplete}} @onShamirSuccess={{action this.onSuccess}} />
    `);
    await click('[data-test-shamir-submit]');
  });

  test('it fetches progress on init when fetchOnInit is true', async function (assert) {
    await render(hbs`
      <ShamirFlow @action="foo" @fetchOnInit={{true}} />
    `);
    assert
      .dom('.shamir-progress')
      .hasText(`${response.progress}/${response.required} keys provided`, 'displays the correct progress');
  });
});
