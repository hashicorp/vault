/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This cannot be called kv-metadata because mirage checks for plural factory names, and metadata and data are considered plural. It will throw an error.
import { Factory, trait } from 'ember-cli-mirage';

// define data outside of factory for linting error: https://github.com/ember-cli/eslint-plugin-ember/issues/202#issuecomment-356255988
const data = {
  path: 'my-secret',
  backend: 'kv-engine',
  cas_required: false,
  created_time: '2018-03-22T02:24:06.945319214Z',
  current_version: 3,
  delete_version_after: '3h25m19s',
  max_versions: 0,
  oldest_version: 0,
  updated_time: '2018-03-22T02:36:43.986212308Z',
  versions: {
    1: {
      created_time: '2018-03-22T02:24:06.945319214Z',
      deletion_time: '',
      destroyed: false,
    },
    2: {
      created_time: '2018-03-22T02:36:33.954880664Z',
      deletion_time: '',
      destroyed: false,
    },
    3: {
      created_time: '2018-03-22T02:36:43.986212308Z',
      deletion_time: '',
      destroyed: false,
    },
  },
};

export default Factory.extend({
  data,

  withCustomMetadata: trait({
    custom_metadata: {
      foo: 'abc',
      bar: '123',
      baz: '5c07d823-3810-48f6-a147-4c06b5219e84',
    },
  }),
});
