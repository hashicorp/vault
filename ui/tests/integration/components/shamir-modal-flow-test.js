import { module, test, skip } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | shamir-modal-flow', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.set('isActive', true);
    this.set('onClose', sinon.spy());
    this.server.get('/sys/replication/dr/secondary/generate-operation-token/attempt', () => {});
  });

  test('it renders with initial content by default', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <ShamirModalFlow
        @action="generate-dr-operation-token"
        @buttonText="Generate token"
        @fetchOnInit=true
        @generateAction=true
        @buttonText="My CTA"
        @onClose={{this.onClose}}
        @isActive={{this.isActive}}
      >
        <p>Inner content goes here</p>
      </ShamirModalFlow>
    `);

    assert
      .dom('[data-test-shamir-modal-body]')
      .hasText('Inner content goes here', 'Template block gets rendered');
    assert.dom('[data-test-shamir-modal-cancel-button]').hasText('Cancel', 'Shows cancel button');
  });

  test('Shows correct content when started', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <ShamirModalFlow
        @started=true
        @action="generate-dr-operation-token"
        @buttonText="Generate token"
        @fetchOnInit=true
        @generateAction=true
        @buttonText="Crazy CTA"
        @onClose={{this.onClose}}
        @isActive={{this.isActive}}
      >
        <p>Inner content goes here</p>
      </ShamirModalFlow>
    `);
    assert.dom('[data-test-shamir-input]').exists('Asks for root key Portion');
    assert.dom('[data-test-shamir-modal-cancel-button]').hasText('Cancel', 'Shows cancel button');
  });

  test('Shows OTP when provided and flow started', async function (assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <ShamirModalFlow
        @encoded_token="my-encoded-token"
        @action="generate-dr-operation-token"
        @buttonText="Generate token"
        @fetchOnInit=true
        @generateAction=true
        @buttonText="Crazy CTA"
        @onClose={{this.onClose}}
        @isActive={{this.isActive}}
      >
        <p>Inner content goes here</p>
      </ShamirModalFlow>
    `);
    assert.dom('[data-test-shamir-encoded-token]').hasText('my-encoded-token', 'Shows encoded token');
    assert.dom('[data-test-shamir-modal-cancel-button]').hasText('Close', 'Shows close button');
  });
  skip('DR Secondary actions', async function () {
    // DR Secondaries cannot be tested yet, but once they can
    // we should add tests for Cancel button functionality
  });
});
