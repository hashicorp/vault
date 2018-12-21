import { later, run } from '@ember/runloop';
import { module, test, skip } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import hbs from 'htmlbars-inline-precompile';

import { create } from 'ember-cli-page-object';
import mountBackendForm from '../../pages/components/mount-backend-form';

import sinon from 'sinon';

const component = create(mountBackendForm);

module('Integration | Component | mount backend form', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
    this.owner.lookup('service:flash-messages').registerTypes(['success', 'danger']);
    this.server = apiStub({ usePassthrough: true });
  });

  hooks.afterEach(function() {
    component.removeContext();
    this.server.shutdown();
  });

  test('it renders', async function(assert) {
    await render(hbs`{{mount-backend-form}}`);
    assert.equal(component.header, 'Enable an authentication method', 'renders auth header in default state');
    assert.ok(component.types.length > 0, 'renders type picker');
  });

  test('it changes path when type is changed', async function(assert) {
    await render(hbs`{{mount-backend-form}}`);
    await component.selectType('aws');
    await component.next();
    assert.equal(component.pathValue, 'aws', 'sets the value of the type');
    await component.back();
    await component.selectType('approle');
    await component.next();
    assert.equal(component.pathValue, 'approle', 'updates the value of the type');
  });

  test('it keeps path value if the user has changed it', async function(assert) {
    await render(hbs`{{mount-backend-form}}`);
    await component.selectType('approle');
    await component.next();
    assert.equal(component.pathValue, 'approle', 'defaults to approle (first in the list)');
    await component.path('newpath');
    await component.back();
    await component.selectType('aws');
    await component.next();
    assert.equal(component.pathValue, 'newpath', 'updates to the value of the type');
  });

  skip('it calls mount success', async function(assert) {
    const spy = sinon.spy();
    this.set('onMountSuccess', spy);
    await render(hbs`{{mount-backend-form onMountSuccess=onMountSuccess}}`);

    await component.mount('approle', 'foo');
    assert.equal(this.server.db.authMethods.length, 1, 'it enables an auth method');
    assert.ok(spy.calledOnce, 'calls the passed success method');
  });

  skip('it calls mount config error', async function(assert) {
    const spy = sinon.spy();
    const spy2 = sinon.spy();
    this.set('onMountSuccess', spy);
    this.set('onConfigError', spy2);
    await render(hbs`{{mount-backend-form onMountSuccess=onMountSuccess onConfigError=onConfigError}}`);

    await component.selectType('kubernetes');
    await component.next().path('bar');
    // kubernetes requires a host + a cert / pem, so only filling the host will error
    await component.fillIn('kubernetesHost', 'host');
    component.submit();
    later(() => run.cancelTimers(), 50);
    await settled();
    assert.equal(this.server.db.authMethods.length, 1, 'it still enables an auth method');
    assert.equal(spy.callCount, 0, 'does not call the success method');
    assert.ok(spy2.calledOnce, 'calls the passed error method');
  });
});
