/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';

export default class RoleSamlModel extends Model {
  @attr('string') ssoServiceURL;
  @attr('string') tokenPollID;
  @attr('string') clientVerifier;
}
