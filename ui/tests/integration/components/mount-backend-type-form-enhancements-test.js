/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | mount-backend/type-form | plugin catalog enhancements', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.setMountType = sinon.spy();
    this.version = this.owner.lookup('service:version');
    this.store = this.owner.lookup('service:store');
  });

  test('it renders basic form elements', async function (assert) {
    await render(hbs`
      <MountBackend::TypeForm 
        @mountCategory="secret" 
        @setMountType={{this.setMountType}} 
      />
    `);

    assert.dom('.field.is-grouped').exists('renders action buttons');
    assert.dom(`${GENERAL.cancelButton}`).exists('shows cancel button');
  });
});
