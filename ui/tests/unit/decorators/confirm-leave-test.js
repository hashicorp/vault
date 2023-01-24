import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import { withConfirmLeave } from 'core/decorators/confirm-leave';
import Route from '@ember/routing/route';

const ControllerStub = {
  get(name) {
    return {
      name,
    };
  },
};
// create class using decorator
// const createClass = (mainModel, otherModels) => {
//   @withConfirmLeave(mainModel, otherModels)
//   class Foo extends Route {
//     @service store;
//     model() {
//       return this.store.createRecord('pki/action');
//     }
//   }
//   return new Foo();
// };

module('Unit | Decorators | ConfirmLeave', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.spy = sinon.spy(console, 'error');
    this.controller = this.owner.lookup('controller:foo');
    this.router = this.owner.lookup('service:router');
    this.store = this.owner.lookup('service:store');
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'abc_example',
        path: 'example/',
        type: 'kubernetes',
      },
    });
  });
  hooks.afterEach(function () {
    this.spy.restore();
  });

  test('it should warn when applying decorator to class that does not extend Model', function (assert) {
    @withConfirmLeave()
    class Foo {} // eslint-disable-line
    const message =
      'withConfirmLeave decorator must be used on instance of ember Route class. Decorator not applied to returned class';
    assert.ok(this.spy.calledWith(message), 'Error is printed to console');
  });

  test('it should work correctly without paths passed', function (assert) {
    // @withConfirmLeave()
    class Foo extends Route {
      model() {
        return this.store.createRecord('pki/action');
      }
      controller = ControllerStub;
    }
    const route = new Foo();
    assert.ok(route, 'route is ok');
    route.setup(this.owner);
    assert.ok(false, 'test');
    // const result = route.willTransition();
    // assert.deepEqual(result, 'hello');
  });
});
