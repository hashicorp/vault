/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ComputedProperty from '@ember/object/computed';
import { Model } from 'vault/app-types';

type CapabilitiesModel = Model & {
  path: string;
  capabilities: Array<string>;
  canCreate: ComputedProperty<boolean | undefined>;
  canDelete: ComputedProperty<boolean | undefined>;
  canList: ComputedProperty<boolean | undefined>;
  canPatch: ComputedProperty<boolean | undefined>;
  canRead: ComputedProperty<boolean | undefined>;
  canSudo: ComputedProperty<boolean | undefined>;
  canUpdate: ComputedProperty<boolean | undefined>;
  // these don't seem to be used anywhere
  // inferring type from key name
  allowedParameters: Array<string>;
  deniedParameters: Array<string>;
};

export default CapabilitiesModel;
export const SUDO_PATHS: string[];
export const SUDO_PATH_PREFIXES: string[];
