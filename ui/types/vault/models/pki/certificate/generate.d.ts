/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { FormField, FormFieldGroups } from 'vault/app-types';
import PkiCertificateBaseModel from './base';

export default class PkiCertificateGenerateModel extends PkiCertificateBaseModel {
  role: string;
  formFields: FormField[];
  formFieldGroups: FormFieldGroups;
}
