/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-issuer-details';

module('Integration | Component | page/pki-issuer-details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.context = { owner: this.engine };
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';
    this.issuer = this.store.createRecord('pki/issuer', { issuerId: 'abcd-efgh' });
  });

  test('it renders with correct toolbar by default', async function (assert) {
    await render(
      hbs`
      <Page::PkiIssuerDetails @issuer={{this.issuer}} />
      <div id="modal-wormhole"></div>
      `,
      this.context
    );

    assert.dom(SELECTORS.rotateRoot).doesNotExist();
    assert.dom(SELECTORS.crossSign).doesNotExist();
    assert.dom(SELECTORS.signIntermediate).doesNotExist();
    assert.dom(SELECTORS.download).hasText('Download');
    assert.dom(SELECTORS.configure).doesNotExist();
  });

  test('it renders toolbar actions depending on passed capabilities', async function (assert) {
    this.set('isRotatable', true);
    this.set('canRotate', true);
    this.set('canCrossSign', true);
    this.set('canSignIntermediate', true);
    this.set('canConfigure', true);

    await render(
      hbs`
      <Page::PkiIssuerDetails
        @issuer={{this.issuer}}
        @isRotatable={{this.isRotatable}}
        @canRotate={{this.canRotate}}
        @canCrossSign={{this.canCrossSign}}
        @canSignIntermediate={{this.canSignIntermediate}}
        @canConfigure={{this.canConfigure}}
      />
      <div id="modal-wormhole"></div>
      `,
      this.context
    );

    assert.dom(SELECTORS.rotateRoot).hasText('Rotate this root');
    assert.dom(SELECTORS.crossSign).hasText('Cross-sign issuers');
    assert.dom(SELECTORS.signIntermediate).hasText('Sign Intermediate');
    assert.dom(SELECTORS.download).hasText('Download');
    assert.dom(SELECTORS.configure).hasText('Configure');

    this.set('canRotate', false);
    this.set('canCrossSign', false);
    this.set('canSignIntermediate', false);
    this.set('canConfigure', false);
    await settled();

    assert.dom(SELECTORS.rotateRoot).doesNotExist();
    assert.dom(SELECTORS.crossSign).doesNotExist();
    assert.dom(SELECTORS.signIntermediate).doesNotExist();
    assert.dom(SELECTORS.download).hasText('Download');
    assert.dom(SELECTORS.configure).doesNotExist();
  });
});
