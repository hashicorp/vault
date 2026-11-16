/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import itemDetails from 'vault/tests/pages/components/identity/item-details';

const component = create(itemDetails);

module('Integration | Component | identity/item details', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.flash = this.owner.lookup('service:flash-messages');
    this.flash.registerTypes(['success', 'danger']);
    this.successSpy = sinon.spy(this.flash, 'success');
    this.dangerSpy = sinon.spy(this.flash, 'danger');
    this.api = this.owner.lookup('service:api');
    this.router = this.owner.lookup('service:router');
    this.refreshStub = sinon.stub(this.router, 'refresh');
    this.entityUpdate = sinon.stub(this.api.identity, 'entityUpdateById').resolves();
  });

  test('it renders the disabled warning and enables the entity through the API', async function (assert) {
    this.model = { id: 'e-1', identityType: 'entity', disabled: true, canEdit: true };
    await render(hbs`<Identity::ItemDetails @model={{this.model}} />`);
    assert.dom('[data-test-disabled-warning]').exists();
    assert.dom('[data-test-disabled-warning]').includesText('This entity is disabled');
    await component.enable();

    assert.true(
      this.entityUpdate.calledOnceWithExactly('e-1', { disabled: false }),
      'only the disabled flag is sent'
    );
    assert.true(
      this.refreshStub.calledOnceWithExactly('vault.cluster.access.identity.entities.show'),
      'only the show route is refreshed so the success flash is not torn down'
    );
    assert.true(this.successSpy.calledOnceWith('Successfully enabled entity: e-1'));
  });

  test('it reports an error and does not refresh when enabling fails', async function (assert) {
    this.entityUpdate.rejects(new Error('permission denied'));
    this.model = { id: 'e-1', identityType: 'entity', disabled: true, canEdit: true };
    await render(hbs`<Identity::ItemDetails @model={{this.model}} />`);
    await component.enable();

    assert.true(this.refreshStub.notCalled);
    assert.true(this.successSpy.notCalled);
    assert.true(this.dangerSpy.calledOnce);
    assert.true(this.dangerSpy.firstCall.args[0].includes('There was a problem enabling the entity: e-1'));
  });

  test('it does not render the button if canEdit is false', async function (assert) {
    this.model = { id: 'e-1', identityType: 'entity', disabled: true };
    await render(hbs`<Identity::ItemDetails @model={{this.model}} />`);
    assert.dom('[data-test-disabled-warning]').exists('shows the warning banner');
    assert.dom('[data-test-enable-identity]').doesNotExist('does not show the enable button');
  });

  test('it does not render the banner when item is enabled', async function (assert) {
    this.model = { id: 'e-1', identityType: 'entity' };
    await render(hbs`<Identity::ItemDetails @model={{this.model}} />`);
    assert.dom('[data-test-disabled-warning]').doesNotExist('does not show the warning banner');
  });

  test('it renders the last updated time', async function (assert) {
    this.model = { id: 'e-1', identityType: 'entity', last_update_time: '2026-01-02T03:04:05Z' };
    await render(hbs`<Identity::ItemDetails @model={{this.model}} />`);
    assert.dom('time[datetime="2026-01-02T03:04:05Z"]').exists();
  });
});
