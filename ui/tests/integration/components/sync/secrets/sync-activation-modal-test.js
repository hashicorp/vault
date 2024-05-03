import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | Secrets::SyncActivationModal', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onClose = () => {};
    this.onError = () => {};

    this.renderComponent = async () => {
      await render(hbs`
      <Secrets::SyncActivationModal @onClose={{this.onClose}} @onError={{this.onError}} />
    `);
    };
  });

  test('it renders with correct text', async function (assert) {
    await this.renderComponent();

    assert.dom(this.element).hasText('Hello');
  });

  test('it disables submit until user has confirmed docs', async function (assert) {
    await this.renderComponent();

    assert.dom(this.element).hasText('Hello');
  });

  module('on submit', function () {
    module('success', function () {
      test('it calls the activate endpoint', function () {});
      test('it transitions back to sync overview', function () {});
    });

    module('on error', function () {
      test('it calls onError', function () {});
      test('it renders an error flash message', function () {});
    });
  });
});
