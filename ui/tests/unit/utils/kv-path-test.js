import { kvDataPath, kvDestroyPath, kvMetadataPath, kvUndeletePath } from 'vault/utils/kv-path';
import { module, test } from 'qunit';

module('Unit | Utility | kv-path utils', function () {
  module('kvDataPath', function () {
    [
      {
        backend: 'basic-backend',
        path: 'secret-path',
        expected: 'basic-backend/data/secret-path',
      },
      {
        backend: 'some/back end',
        path: 'my/secret/path',
        expected: 'some/back%20end/data/my/secret/path',
      },
      {
        backend: 'some/back end',
        path: 'my/secret/path',
        version: 3,
        expected: 'some/back%20end/data/my/secret/path?version=3',
      },
    ].forEach((t, idx) => {
      test(`kvDataPath ${idx}`, function (assert) {
        const result = kvDataPath(t.backend, t.path, t.version);
        assert.strictEqual(result, t.expected);
      });
    });
  });

  module('kvMetadataPath', function () {
    [
      {
        backend: 'basic-backend',
        path: 'secret-path',
        expected: 'basic-backend/metadata/secret-path',
      },
      {
        backend: 'some/back end',
        path: 'my/secret/path',
        expected: 'some/back%20end/metadata/my/secret/path',
      },
    ].forEach((t, idx) => {
      test(`kvMetadataPath ${idx}`, function (assert) {
        const result = kvMetadataPath(t.backend, t.path);
        assert.strictEqual(result, t.expected);
      });
    });
  });

  module('kvDestroyPath', function () {
    [
      {
        backend: 'basic-backend',
        path: 'secret-path',
        expected: 'basic-backend/destroy/secret-path',
      },
      {
        backend: 'some/back end',
        path: 'my/secret/path',
        expected: 'some/back%20end/destroy/my/secret/path',
      },
    ].forEach((t, idx) => {
      test(`kvDestroyPath ${idx}`, function (assert) {
        const result = kvDestroyPath(t.backend, t.path);
        assert.strictEqual(result, t.expected);
      });
    });
  });

  module('kvUndeletePath', function () {
    [
      {
        backend: 'basic-backend',
        path: 'secret-path',
        expected: 'basic-backend/undelete/secret-path',
      },
      {
        backend: 'some/back end',
        path: 'my/secret/path',
        expected: 'some/back%20end/undelete/my/secret/path',
      },
    ].forEach((t, idx) => {
      test(`kvUndeletePath ${idx}`, function (assert) {
        const result = kvUndeletePath(t.backend, t.path);
        assert.strictEqual(result, t.expected);
      });
    });
  });
});
