/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ComputedProperty from '@ember/object/computed';
import Model from '@ember-data/model';

interface CapabilitiesModel extends Model {
  path: string;
  capabilities: Array<string>;
  canSudo: ComputedProperty<boolean | undefined>;
  canRead: ComputedProperty<boolean | undefined>;
  canCreate: ComputedProperty<boolean | undefined>;
  canUpdate: ComputedProperty<boolean | undefined>;
  canDelete: ComputedProperty<boolean | undefined>;
  canList: ComputedProperty<boolean | undefined>;
  // these don't seem to be used anywhere
  // inferring type from key name
  allowedParameters: Array<string>;
  deniedParameters: Array<string>;
}

export default CapabilitiesModel;
export const SUDO_PATHS: string[];
export const SUDO_PATH_PREFIXES: string[];
