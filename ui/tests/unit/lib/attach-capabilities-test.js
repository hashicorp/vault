/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import attachCapabilities from 'vault/lib/attach-capabilities';
import apiPath from 'vault/utils/api-path';

const MODEL_TYPE = 'test-form-model';

const makeModelClass = () => {
  return Model.extend();
};

module('Unit | lib | attach capabilities', function (hooks) {
  setupTest(hooks);

  test('it attaches passed capabilities', function (assert) {
    let mc = makeModelClass();
    mc = attachCapabilities(mc, {
      updatePath: apiPath`update/{'id'}`,
      deletePath: apiPath`delete/{'id'}`,
    });
    let relationship = mc.relationshipsByName.get('updatePath');

    assert.strictEqual(relationship.key, 'updatePath', 'has updatePath relationship');
    assert.strictEqual(relationship.kind, 'belongsTo', 'kind of relationship is belongsTo');
    assert.strictEqual(relationship.type, 'capabilities', 'updatePath is a related capabilities model');

    relationship = mc.relationshipsByName.get('deletePath');
    assert.strictEqual(relationship.key, 'deletePath', 'has deletePath relationship');
    assert.strictEqual(relationship.kind, 'belongsTo', 'kind of relationship is belongsTo');
    assert.strictEqual(relationship.type, 'capabilities', 'deletePath is a related capabilities model');
  });

  test('it adds a static method to the model class', function (assert) {
    let mc = makeModelClass();
    mc = attachCapabilities(mc, {
      updatePath: apiPath`update/{'id'}`,
      deletePath: apiPath`delete/{'id'}`,
    });
    const relatedCapabilities = !!mc.relatedCapabilities && typeof mc.relatedCapabilities === 'function';
    assert.true(relatedCapabilities, 'model class now has a relatedCapabilities static function');
  });

  test('calling static method with single response JSON-API document adds expected relationships', function (assert) {
    let mc = makeModelClass();
    mc = attachCapabilities(mc, {
      updatePath: apiPath`update/${'id'}`,
      deletePath: apiPath`delete/${'id'}`,
    });
    const jsonAPIDocSingle = {
      data: {
        id: 'test',
        type: MODEL_TYPE,
        attributes: {},
        relationships: {},
      },
      included: [],
    };

    const expected = {
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

    assert.strictEqual(
      Object.keys(jsonAPIDocSingle.data.relationships).length,
      2,
      'document now has 2 relationships'
    );
    assert.deepEqual(jsonAPIDocSingle, expected, 'has the exected new document structure');
  });

  test('calling static method with an arrary response JSON-API document adds expected relationships', function (assert) {
    let mc = makeModelClass();
    mc = attachCapabilities(mc, {
      updatePath: apiPath`update/${'id'}`,
      deletePath: apiPath`delete/${'id'}`,
    });
    const jsonAPIDocSingle = {
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

    const expected = {
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
