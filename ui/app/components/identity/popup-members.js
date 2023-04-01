/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import Base from './_popup-base';

export default Base.extend({
  model: alias('params.firstObject'),

  groupArray: computed('params', function () {
    return this.params.objectAt(1);
  }),

  memberId: computed('params', function () {
    return this.params.objectAt(2);
  }),

  messageArgs(/*model, groupArray, memberId*/) {
    return [...arguments];
  },

  successMessage(model, groupArray, memberId) {
    return `Successfully removed '${memberId}' from the group`;
  },

  errorMessage(e, model, groupArray, memberId) {
    const error = e.errors ? e.errors.join(' ') : e.message;
    return `There was a problem removing '${memberId}' from the group - ${error}`;
  },

  transaction(model, groupArray, memberId) {
    const members = model.get(groupArray);
    model.set(groupArray, members.without(memberId));
    return model.save();
  },
});
