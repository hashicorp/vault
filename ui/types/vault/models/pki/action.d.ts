/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model from '@ember-data/model';
import { FormField, ModelValidations } from 'vault/app-types';
import CapabilitiesModel from '../capabilities';

export default class PkiActionModel extends Model {
  secretMountPath: unknown;
  pemBundle: string;
  type: string;
  actionType: string | null;
  get backend(): string;
  // apiPaths for capabilities
  importBundlePath: Promise<CapabilitiesModel>;
  generateIssuerRootPath: Promise<CapabilitiesModel>;
  generateIssuerCsrPath: Promise<CapabilitiesModel>;
  crossSignPath: string;
  allFields: Array<FormField>;
  validate(): ModelValidations;
  // Capabilities
  get canImportBundle(): boolean;
  get canGenerateIssuerRoot(): boolean;
  get canGenerateIssuerIntermediate(): boolean;
  get canCrossSign(): boolean;
}
