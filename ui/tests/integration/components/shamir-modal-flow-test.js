import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | shamir-modal-flow', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('isActive', true);
    this.set('onClose', sinon.spy());
  });

  test('it renders with initial content by default', async function(assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <ShamirModalFlow
        @action="generate-dr-operation-token"
        @buttonText="Generate token"
        @fetchOnInit=true
        @generateAction=true
        @buttonText="My CTA"
        @onClose={{onClose}}
        @isActive={{isActive}}
      >
        <p>Inner content goes here</p>
      </ShamirModalFlow>
    `);

    assert
      .dom('[data-test-shamir-modal-body]')
      .hasText('Inner content goes here', 'Template block gets rendered');
    assert.dom('[data-test-shamir-modal-cancel-button]').hasText('Cancel', 'Shows cancel button');
  });

  test('Shows correct content when started', async function(assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <ShamirModalFlow
        @started=true
        @action="generate-dr-operation-token"
        @buttonText="Generate token"
        @fetchOnInit=true
        @generateAction=true
        @buttonText="Crazy CTA"
        @onClose={{onClose}}
        @isActive={{isActive}}
      >
        <p>Inner content goes here</p>
      </ShamirModalFlow>
    `);
    assert.dom('[data-test-shamir-input]').exists('Asks for Master Key Portion');
    assert.dom('[data-test-shamir-modal-cancel-button]').hasText('Cancel', 'Shows cancel button');
  });

  test('Shows OTP when provided and flow started', async function(assert) {
    await render(hbs`
      <div id="modal-wormhole"></div>
      <ShamirModalFlow
        @encoded_token="my-encoded-token"
        @action="generate-dr-operation-token"
        @buttonText="Generate token"
        @fetchOnInit=true
        @generateAction=true
        @buttonText="Crazy CTA"
        @onClose={{onClose}}
        @isActive={{isActive}}
      >
        <p>Inner content goes here</p>
      </ShamirModalFlow>
    `);
    assert.dom('[data-test-shamir-encoded-token]').hasText('my-encoded-token', 'Shows encoded token');
    assert.dom('[data-test-shamir-modal-cancel-button]').hasText('Close', 'Shows close button');
  });
  /*
  test('DR Secondary actions', async function (assert) {
    // DR Secondaries cannot be tested yet, but once they can
    // we should add tests for Cancel button functionality
  })
  */
});
