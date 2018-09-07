import { moduleFor, test } from 'ember-qunit';
import { STORAGE_KEYS, DEFAULTS } from 'vault/helpers/wizard-constants';

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
    args: [STORAGE_KEYS.TUTORIAL_STATE],
    expectedResults: {
      storage: [{ key: STORAGE_KEYS.TUTORIAL_STATE, value: 'idle' }],
    },
  },
  {
    method: 'saveExtState',
    args: [STORAGE_KEYS.TUTORIAL_STATE, 'test'],
    expectedResults: {
      storage: [{ key: STORAGE_KEYS.TUTORIAL_STATE, value: 'test' }],
    },
  },
  {
    method: 'storageHasKey',
    args: ['fake-key'],
    expectedResults: { value: false },
  },
  {
    method: 'storageHasKey',
    args: [STORAGE_KEYS.TUTORIAL_STATE],
    expectedResults: { value: true },
  },
  {
    method: 'handleDismissed',
    args: [],
    expectedResults: {
      storage: [
        { key: STORAGE_KEYS.FEATURE_STATE, value: undefined },
        { key: STORAGE_KEYS.FEATURE_LIST, value: undefined },
        { key: STORAGE_KEYS.COMPONENT_STATE, value: undefined },
      ],
    },
  },
  {
    method: 'handlePaused',
    args: [],
    properties: {
      expectedURL: 'this/is/a/url',
      expectedRouteName: 'this.is.a.route',
    },
    expectedResults: {
      storage: [
        { key: STORAGE_KEYS.RESUME_URL, value: 'this/is/a/url' },
        { key: STORAGE_KEYS.RESUME_ROUTE, value: 'this.is.a.route' },
      ],
    },
  },
  {
    method: 'handlePaused',
    args: [],
    expectedResults: {
      storage: [
        { key: STORAGE_KEYS.RESUME_URL, value: undefined },
        { key: STORAGE_KEYS.RESUME_ROUTE, value: undefined },
      ],
    },
  },
  {
    method: 'restartGuide',
    args: [],
    expectedResults: {
      props: [
        { prop: 'currentState', value: 'active.select' },
        { prop: 'featureComponent', value: 'wizard/features-selection' },
        { prop: 'tutorialComponent', value: 'wizard/tutorial-active' },
      ],
      storage: [
        { key: STORAGE_KEYS.FEATURE_STATE, value: undefined },
        { key: STORAGE_KEYS.FEATURE_LIST, value: undefined },
        { key: STORAGE_KEYS.COMPONENT_STATE, value: undefined },
        { key: STORAGE_KEYS.TUTORIAL_STATE, value: 'active.select' },
        { key: STORAGE_KEYS.COMPLETED_FEATURES, value: undefined },
        { key: STORAGE_KEYS.RESUME_URL, value: undefined },
        { key: STORAGE_KEYS.RESUME_ROUTE, value: undefined },
      ],
    },
  },
  {
    method: 'saveState',
    args: [
      'currentState',
      {
        value: {
          init: {
            active: 'login',
          },
        },
        actions: [{ type: 'render', level: 'feature', component: 'wizard/init-login' }],
      },
    ],
    expectedResults: {
      props: [{ prop: 'currentState', value: 'init.active.login' }],
    },
  },
  {
    method: 'saveState',
    args: [
      'currentState',
      {
        value: {
          active: 'login',
        },
        actions: [{ type: 'render', level: 'feature', component: 'wizard/init-login' }],
      },
    ],
    expectedResults: {
      props: [{ prop: 'currentState', value: 'active.login' }],
    },
  },
  {
    method: 'saveState',
    args: ['currentState', 'login'],
    expectedResults: {
      props: [{ prop: 'currentState', value: 'login' }],
    },
  },
];

testCases.forEach(testCase => {
  let store = storage();
  test(`${testCase.method}`, function(assert) {
    let wizard = this.subject({
      storage() {
        return store;
      },
    });

    if (testCase.properties) {
      wizard.setProperties(testCase.properties);
    } else {
      wizard.setProperties(DEFAULTS);
    }

    if (testCase.storage) {
      testCase.storage.forEach(item => wizard.storage().setItem(item.key, item.value));
    }

    let result = wizard[testCase.method](...testCase.args);
    if (testCase.expectedResults.props) {
      testCase.expectedResults.props.forEach(property => {
        assert.deepEqual(
          wizard.get(property.prop),
          property.value,
          `${testCase.method} creates correct value for ${property.prop}`
        );
      });
    }
    if (testCase.expectedResults.storage) {
      testCase.expectedResults.storage.forEach(item => {
        assert.deepEqual(
          wizard.storage().getItem(item.key),
          item.value,
          `${testCase.method} creates correct storage state for ${item.key}`
        );
      });
    }
    if (testCase.expectedResults.value !== null && testCase.expectedResults.value !== undefined) {
      assert.equal(result, testCase.expectedResults.value, `${testCase.method} gives correct value`);
    }
  });
});
