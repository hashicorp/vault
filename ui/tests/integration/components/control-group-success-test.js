import { later, run } from '@ember/runloop';
import { resolve } from 'rsvp';
import Service from '@ember/service';
import { moduleForComponent, test } from 'ember-qunit';
import wait from 'ember-test-helpers/wait';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import controlGroupSuccess from '../../pages/components/control-group-success';

const component = create(controlGroupSuccess);

const controlGroupService = Service.extend({
  deleteControlGroupToken: sinon.stub(),
  markTokenForUnwrap: sinon.stub(),
});

const routerService = Service.extend({
  transitionTo: sinon.stub().returns(resolve()),
});

const storeService = Service.extend({
  adapterFor() {
    return {
      toolAction() {
        return resolve({ data: { foo: 'bar' } });
      },
    };
  },
});

moduleForComponent('control-group-success', 'Integration | Component | control group success', {
  integration: true,
  beforeEach() {
    component.setContext(this);
    this.register('service:control-group', controlGroupService);
    this.inject.service('controlGroup');
    this.register('service:router', routerService);
    this.inject.service('router');
    this.register('service:store', storeService);
    this.inject.service('store');
  },

  afterEach() {
    component.removeContext();
  },
});

const MODEL = {
  approved: false,
  requestPath: 'foo/bar',
  id: 'accessor',
  requestEntity: { id: 'requestor', name: 'entity8509' },
  reload: sinon.stub(),
};
test('render with saved token', function(assert) {
  let response = {
    uiParams: { url: '/foo' },
    token: 'token',
  };
  this.set('model', MODEL);
  this.set('response', response);
  this.render(hbs`{{control-group-success model=model controlGroupResponse=response }}`);
  assert.ok(component.showsNavigateMessage, 'shows unwrap message');
  component.navigate();
  later(() => run.cancelTimers(), 50);
  return wait().then(() => {
    assert.ok(this.get('controlGroup').markTokenForUnwrap.calledOnce, 'marks token for unwrap');
    assert.ok(this.get('router').transitionTo.calledOnce, 'calls router transition');
  });
});

test('render without token', function(assert) {
  this.set('model', MODEL);
  this.render(hbs`{{control-group-success model=model}}`);
  assert.ok(component.showsUnwrapForm, 'shows unwrap form');
  component.token('token');
  component.unwrap();

  later(() => run.cancelTimers(), 50);
  return wait().then(() => {
    assert.ok(component.showsJsonViewer, 'shows unwrapped data');
  });
});
