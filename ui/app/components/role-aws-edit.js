/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isBlank } from '@ember/utils';
import { set } from '@ember/object';
import RoleEdit from './role-edit';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default RoleEdit.extend({
  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      // all of the attributes with fieldValue:'id' are called `name`
      const modelId = this.model.id || this.model.name;
      // prevent from submitting if there's no key
      // maybe do something fancier later
      if (type === 'create' && isBlank(modelId)) {
        return;
      }
      var credential_type = this.model.credential_type;
      if (credential_type == 'iam_user') {
        set(this, 'model.role_arns', []);
      }
      if (credential_type == 'assumed_role') {
        set(this, 'model.policy_arns', []);
      }
      if (credential_type == 'federation_token') {
        set(this, 'model.role_arns', []);
        set(this, 'model.policy_arns', []);
      }

      var policy_document = this.model.policy_document;
      if (policy_document == '{}') {
        set(this, 'model.policy_document', '');
      }

      this.persist('save', () => {
        this.hasDataChanges();
        this.transitionToRoute(SHOW_ROUTE, modelId);
      });
    },

    codemirrorUpdated(attr, val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;

      if (!hasErrors) {
        set(this.model, attr, val);
      }
    },
  },
});
