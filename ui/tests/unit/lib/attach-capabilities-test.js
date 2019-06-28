import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import attachCapabilities from 'vault/lib/attach-capabilities';
import apiPath from 'vault/utils/api-path';
import { get } from '@ember/object';

let MODEL_TYPE = 'test-form-model';

module('Unit | lib | attach capabilities', function(hooks) {
  setupTest(hooks);

  test('it attaches passed capabilities', function(assert) {
    let mc = this.owner.lookup('service:store').modelFor(MODEL_TYPE);
    mc = attachCapabilities(mc, {
      updatePath: apiPath`update/{'id'}`,
      deletePath: apiPath`delete/{'id'}`,
    });
    let relationship = get(mc, 'relationshipsByName').get('updatePath');

    assert.equal(relationship.key, 'updatePath', 'has updatePath relationship');
    assert.equal(relationship.kind, 'belongsTo', 'kind of relationship is belongsTo');
    assert.equal(relationship.type, 'capabilities', 'updatePath is a related capabilities model');

    relationship = get(mc, 'relationshipsByName').get('deletePath');
    assert.equal(relationship.key, 'deletePath', 'has deletePath relationship');
    assert.equal(relationship.kind, 'belongsTo', 'kind of relationship is belongsTo');
    assert.equal(relationship.type, 'capabilities', 'deletePath is a related capabilities model');
  });

  test('it adds a static method to the model class', function(assert) {
    let mc = this.owner.lookup('service:store').modelFor(MODEL_TYPE);
    mc = attachCapabilities(mc, {
      updatePath: apiPath`update/{'id'}`,
      deletePath: apiPath`delete/{'id'}`,
    });
    assert.ok(
      !!mc.relatedCapabilities && typeof mc.relatedCapabilities === 'function',
      'model class now has a relatedCapabilities static function'
    );
  });

  test('calling static method with single response JSON-API document adds expected relationships', function(assert) {
    let mc = this.owner.lookup('service:store').modelFor(MODEL_TYPE);
    mc = attachCapabilities(mc, {
      updatePath: apiPath`update/${'id'}`,
      deletePath: apiPath`delete/${'id'}`,
    });
    let jsonAPIDocSingle = {
      data: {
        id: 'test',
        type: MODEL_TYPE,
        attributes: {},
        relationships: {},
      },
      included: [],
    };

    let expected = {
      data: {
        id: 'test',
        type: MODEL_TYPE,
        attributes: {},
        relationships: {
          updatePath: {
            data: {
              type: 'capabilities',
              id: 'update/test',
            },
          },
          deletePath: {
            data: {
              type: 'capabilities',
              id: 'delete/test',
            },
          },
        },
      },
      included: [],
    };

    mc.relatedCapabilities(jsonAPIDocSingle);

    assert.equal(
      Object.keys(jsonAPIDocSingle.data.relationships).length,
      2,
      'document now has 2 relationships'
    );
    assert.deepEqual(jsonAPIDocSingle, expected, 'has the exected new document structure');
  });

  test('calling static method with an arrary response JSON-API document adds expected relationships', function(assert) {
    let mc = this.owner.lookup('service:store').modelFor(MODEL_TYPE);
    mc = attachCapabilities(mc, {
      updatePath: apiPath`update/${'id'}`,
      deletePath: apiPath`delete/${'id'}`,
    });
    let jsonAPIDocSingle = {
      data: [
        {
          id: 'test',
          type: MODEL_TYPE,
          attributes: {},
          relationships: {},
        },
        {
          id: 'foo',
          type: MODEL_TYPE,
          attributes: {},
          relationships: {},
        },
      ],
      included: [],
    };

    let expected = {
      data: [
        {
          id: 'test',
          type: MODEL_TYPE,
          attributes: {},
          relationships: {
            updatePath: {
              data: {
                type: 'capabilities',
                id: 'update/test',
              },
            },
            deletePath: {
              data: {
                type: 'capabilities',
                id: 'delete/test',
              },
            },
          },
        },
        {
          id: 'foo',
          type: MODEL_TYPE,
          attributes: {},
          relationships: {
            updatePath: {
              data: {
                type: 'capabilities',
                id: 'update/foo',
              },
            },
            deletePath: {
              data: {
                type: 'capabilities',
                id: 'delete/foo',
              },
            },
          },
        },
      ],
      included: [],
    };
    mc.relatedCapabilities(jsonAPIDocSingle);
    assert.deepEqual(jsonAPIDocSingle, expected, 'has the exected new document structure');
  });
});
