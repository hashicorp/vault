/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { waitUntil } from '@ember/test-helpers';
import Ember from 'ember';

export default function waitForError(opts) {
  const orig = Ember.onerror;

  let error = null;
  Ember.onerror = (err) => {
    error = err;
  };

  return waitUntil(() => error, opts).finally(() => {
    Ember.onerror = orig;
  });
}
