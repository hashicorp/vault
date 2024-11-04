/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import LdapRoleModel from '../role';
import { attr } from '@ember-data/model';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';

const validations = {
  username: [
    {
      validator: (model) => (!model.username ? false : true),
      message: 'Username is required.',
    },
  ],
  rotation_period: [
    {
      validator: (model) => (!model.rotation_period ? false : true),
      message: 'Rotation Period is required.',
    },
  ],
};

// determines form input rendering order
@withFormFields(['name', 'username', 'dn', 'rotation_period'])
@withModelValidations(validations)
export default class LdapRoleStaticModel extends LdapRoleModel {
  type = 'static';
  roleUri = 'static-role';
  credsUri = 'static-cred';

  @attr('string', {
    label: 'Distinguished name',
    subText: 'Distinguished name (DN) of entry Vault should manage.',
  })
  dn;

  @attr('string', {
    label: 'Username',
    subText:
      "The name of the user to be used when logging in. This is useful when DN isn't used for login purposes.",
  })
  username;

  @attr({
    editType: 'ttl',
    label: 'Rotation period',
    helperTextEnabled:
      'Specifies the amount of time Vault should wait before rotating the password. The minimum is 5 seconds.',
    hideToggle: true,
  })
  rotation_period;

  get canRotateStaticCreds() {
    return this.staticRotateCredsPath.get('canCreate') !== false;
  }
}
