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
  cas_required: true,
  created_time: '2018-03-22T02:24:06.945319214Z',
  current_version: 4,
  delete_version_after: '3h25m19s',
  max_versions: 15,
  oldest_version: 0,
  updated_time: '2018-03-22T02:36:43.986212308Z',
  versions: {
    1: {
      created_time: '2023-07-20T02:12:09.11529Z',
      deletion_time: '',
      destroyed: false,
    },
    2: {
      created_time: '2023-07-20T02:15:35.86465Z',
      deletion_time: '2023-07-25T00:36:19.950545Z',
      destroyed: false,
    },
    3: {
      created_time: '2023-07-20T02:15:40.164549Z',
      deletion_time: '',
      destroyed: true,
    },
    4: {
      created_time: '2023-07-21T03:11:58.095971Z',
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

  withCustomPath: trait({
    path(i) {
      return `my-secret-${i}`;
    },
  }),
});
