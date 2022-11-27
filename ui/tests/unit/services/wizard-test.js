/* eslint qunit/no-conditional-assertions: "warn" */
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import { STORAGE_KEYS, DEFAULTS } from 'vault/helpers/wizard-constants';

module('Unit | Service | wizard', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.router.reopen({
      transitionTo: sinon.stub().returns({
        followRedirects: function () {
          return {
            then: function (callback) {
              callback();
            },
          };
        },
      }),
      urlFor: sinon.stub().returns('/ui/vault/foo'),
    });
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

  const testCases = [
    {
      method: 'getExtState',
      args: [STORAGE_KEYS.TUTORIAL_STATE],
      expectedResults: {
        storage: [{ key: STORAGE_KEYS.TUTORIAL_STATE, value: 'idle' }],
      },
      assertCount: 1,
    },
    {
      method: 'saveExtState',
      args: [STORAGE_KEYS.TUTORIAL_STATE, 'test'],
      expectedResults: {
        storage: [{ key: STORAGE_KEYS.TUTORIAL_STATE, value: 'test' }],
      },
      assertCount: 1,
    },
    {
      method: 'storageHasKey',
      args: ['fake-key'],
      expectedResults: { value: false },
      assertCount: 1,
    },
    {
      method: 'storageHasKey',
      args: [STORAGE_KEYS.TUTORIAL_STATE],
      expectedResults: { value: true },
      assertCount: 1,
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
      assertCount: 3,
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
      assertCount: 2,
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
      assertCount: 2,
    },
    {
      method: 'handleResume',
      storage: [
        { key: STORAGE_KEYS.RESUME_URL, value: 'this/is/a/url' },
        { key: STORAGE_KEYS.RESUME_ROUTE, value: 'this.is.a.route' },
      ],
      args: [],
      expectedResults: {
        props: [
          { prop: 'expectedURL', value: 'this/is/a/url' },
          { prop: 'expectedRouteName', value: 'this.is.a.route' },
        ],
        storage: [
          { key: STORAGE_KEYS.RESUME_URL, value: undefined },
          { key: STORAGE_KEYS.RESUME_ROUTE, value: 'this.is.a.route' },
        ],
      },
      assertCount: 4,
    },
    {
      method: 'handleResume',
      args: [],
      expectedResults: {
        storage: [
          { key: STORAGE_KEYS.RESUME_URL, value: undefined },
          { key: STORAGE_KEYS.RESUME_ROUTE, value: undefined },
        ],
      },
      assertCount: 2,
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
          { key: STORAGE_KEYS.FEATURE_STATE_HISTORY, value: undefined },
          { key: STORAGE_KEYS.FEATURE_LIST, value: undefined },
          { key: STORAGE_KEYS.COMPONENT_STATE, value: undefined },
          { key: STORAGE_KEYS.TUTORIAL_STATE, value: 'active.select' },
          { key: STORAGE_KEYS.COMPLETED_FEATURES, value: undefined },
          { key: STORAGE_KEYS.RESUME_URL, value: undefined },
          { key: STORAGE_KEYS.RESUME_ROUTE, value: undefined },
        ],
      },
      assertCount: 11,
    },
    {
      method: 'clearFeatureData',
      args: [],
      expectedResults: {
        props: [
          { prop: 'currentMachine', value: null },
          { prop: 'featureMachineHistory', value: null },
        ],
        storage: [
          { key: STORAGE_KEYS.FEATURE_STATE, value: undefined },
          { key: STORAGE_KEYS.FEATURE_STATE_HISTORY, value: undefined },
          { key: STORAGE_KEYS.FEATURE_LIST, value: undefined },
          { key: STORAGE_KEYS.COMPONENT_STATE, value: undefined },
        ],
      },
      assertCount: 6,
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
      assertCount: 1,
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
      assertCount: 1,
    },
    {
      method: 'saveState',
      args: ['currentState', 'login'],
      expectedResults: {
        props: [{ prop: 'currentState', value: 'login' }],
      },
      assertCount: 1,
    },
    {
      method: 'saveFeatureHistory',
      args: ['idle'],
      properties: { featureList: ['policies', 'tools'] },
      storage: [{ key: STORAGE_KEYS.COMPLETED_FEATURES, value: ['secrets'] }],
      expectedResults: {
        props: [{ prop: 'featureMachineHistory', value: null }],
      },
      assertCount: 1,
    },
    {
      method: 'saveFeatureHistory',
      args: ['idle'],
      properties: { featureList: ['policies', 'tools'] },
      storage: [],
      expectedResults: {
        props: [{ prop: 'featureMachineHistory', value: ['idle'] }],
      },
      assertCount: 1,
    },
    {
      method: 'saveFeatureHistory',
      args: ['idle'],
      properties: { featureList: ['policies', 'tools'] },
      storage: [],
      expectedResults: {
        props: [{ prop: 'featureMachineHistory', value: ['idle'] }],
      },
      assertCount: 1,
    },
    {
      method: 'saveFeatureHistory',
      args: ['idle'],
      properties: { featureMachineHistory: [], featureList: ['policies', 'tools'] },
      storage: [{ key: STORAGE_KEYS.COMPLETED_FEATURES, value: ['secrets'] }],
      expectedResults: {
        props: [{ prop: 'featureMachineHistory', value: ['idle'] }],
        storage: [{ key: STORAGE_KEYS.FEATURE_STATE_HISTORY, value: ['idle'] }],
      },
      assertCount: 2,
    },
    {
      method: 'saveFeatureHistory',
      args: ['idle'],
      properties: { featureMachineHistory: null, featureList: ['policies', 'tools'] },
      storage: [{ key: STORAGE_KEYS.COMPLETED_FEATURES, value: ['secrets'] }],
      expectedResults: {
        props: [{ prop: 'featureMachineHistory', value: null }],
      },
      assertCount: 1,
    },
    {
      method: 'saveFeatureHistory',
      args: ['create'],
      properties: { featureMachineHistory: ['idle'], featureList: ['policies', 'tools'] },
      storage: [{ key: STORAGE_KEYS.COMPLETED_FEATURES, value: ['secrets'] }],
      expectedResults: {
        props: [{ prop: 'featureMachineHistory', value: ['idle', 'create'] }],
        storage: [{ key: STORAGE_KEYS.FEATURE_STATE_HISTORY, value: ['idle', 'create'] }],
      },
      assertCount: 2,
    },
    {
      method: 'saveFeatureHistory',
      args: ['create'],
      properties: { featureMachineHistory: ['idle'], featureList: ['policies', 'tools'] },
      storage: [
        { key: STORAGE_KEYS.COMPLETED_FEATURES, value: ['secrets'] },
        { key: STORAGE_KEYS.FEATURE_STATE_HISTORY, value: ['idle', 'create'] },
      ],
      expectedResults: {
        props: [{ prop: 'featureMachineHistory', value: ['idle', 'create'] }],
        storage: [{ key: STORAGE_KEYS.FEATURE_STATE_HISTORY, value: ['idle', 'create'] }],
      },
      assertCount: 2,
    },
    {
      method: 'startFeature',
      args: [],
      properties: { featureList: ['secrets', 'tools'] },
      expectedResults: {
        props: [
          { prop: 'featureState', value: 'idle' },
          { prop: 'currentMachine', value: 'secrets' },
        ],
      },
      assertCount: 2,
    },
    {
      method: 'saveFeatures',
      args: [['secrets', 'tools']],
      expectedResults: {
        props: [{ prop: 'featureList', value: ['secrets', 'tools'] }],
        storage: [{ key: STORAGE_KEYS.FEATURE_LIST, value: ['secrets', 'tools'] }],
      },
      assertCount: 2,
    },
  ];

  testCases.forEach((testCase) => {
    const store = storage();
    test(`${testCase.method}`, function (assert) {
      assert.expect(testCase.assertCount);
      const wizard = this.owner.factoryFor('service:wizard').create({
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
        testCase.storage.forEach((item) => wizard.storage().setItem(item.key, item.value));
      }

      const result = wizard[testCase.method](...testCase.args);
      if (testCase.expectedResults.props) {
        testCase.expectedResults.props.forEach((property) => {
          assert.deepEqual(
            wizard.get(property.prop),
            property.value,
            `${testCase.method} creates correct value for ${property.prop}`
          );
        });
      }
      if (testCase.expectedResults.storage) {
        testCase.expectedResults.storage.forEach((item) => {
          assert.deepEqual(
            wizard.storage().getItem(item.key),
            item.value,
            `${testCase.method} creates correct storage state for ${item.key}`
          );
        });
      }
      if (testCase.expectedResults.value !== null && testCase.expectedResults.value !== undefined) {
        assert.strictEqual(result, testCase.expectedResults.value, `${testCase.method} gives correct value`);
      }
    });
  });
});
