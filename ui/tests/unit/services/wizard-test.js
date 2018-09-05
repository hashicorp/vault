import { moduleFor, test } from 'ember-qunit';
//import { saveExtState, getExtState, storageHasKey } from 'vault/services/wizard';
//import sinon from 'sinon';

const TUTORIAL_STATE = 'vault:ui-tutorial-state';
const FEATURE_LIST = 'vault:ui-feature-list';
const FEATURE_STATE = 'vault:ui-feature-state';
// const COMPLETED_FEATURES = 'vault:ui-completed-list';
const COMPONENT_STATE = 'vault:ui-component-state';
// const RESUME_URL = 'vault:ui-tutorial-resume-url';
// const RESUME_ROUTE = 'vault:ui-tutorial-resume-route';

moduleFor('service:wizard', 'Unit | Service | wizard', {
  beforeEach() {},
  afterEach() {},
});

function storage() {
  return {
    items: {},
    getItem(key) {
      var item = this.items[key];
      return item && JSON.parse(item);
    },

    setItem(key, val) {
      return (this.items[key] = JSON.stringify(val));
    },

    removeItem(key) {
      delete this.items[key];
    },

    keys() {
      return Object.keys(this.items);
    },
  };
}

let testCases = [
  {
    method: 'getExtState',
    args: [TUTORIAL_STATE],
    expectedResults: {
      storage: [{ key: TUTORIAL_STATE, value: 'idle' }],
    },
  },
  {
    method: 'saveExtState',
    args: [TUTORIAL_STATE, 'test'],
    expectedResults: {
      storage: [{ key: TUTORIAL_STATE, value: 'test' }],
    },
  },
  {
    method: 'storageHasKey',
    args: [TUTORIAL_STATE],
    expectedResults: { value: false },
  },
  {
    method: 'handleDismissed',
    args: [],
    expectedResults: {
      storage: [
        { key: FEATURE_STATE, value: undefined },
        { key: FEATURE_LIST, value: undefined },
        { key: COMPONENT_STATE, value: undefined },
      ],
    },
  },
];

test('it handles localStorage properly', function(assert) {
  let store = storage();
  let wizard = this.subject({
    storage() {
      return store;
    },
  });

  testCases.forEach(testCase => {
    let result = wizard[testCase.method](...testCase.args);
    if (testCase.expectedResults.props) {
      testCase.expectedResults.props.forEach(property => {
        assert.equal(wizard.get(property.prop), property.value, `${testCase.method} creates correct state`);
      });
    }
    if (testCase.expectedResults.storage) {
      testCase.expectedResults.storage.forEach(item => {
        assert.equal(
          wizard.storage().getItem(item.key),
          item.value,
          `${testCase.method} creates correct storage state`
        );
      });
    }
    if (testCase.expectedResults.value) {
      assert.equal(result, testCase.expectedResults.value, `${testCase.method} gives correct value`);
    }
  });
});
