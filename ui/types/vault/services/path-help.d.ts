/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';

import type { PathInfo } from 'vault/utils/openapi-helpers';
import type FormField from 'vault/utils/forms/field';
import type FormFieldGroup from 'vault/utils/forms/field-group';

export interface FormMetadata {
  formFields: FormField[];
  formFieldGroups: FormFieldGroup[];
  // name of the path parameter field (e.g. `username`), sent in the URL rather than the body
  identifier?: string;
}

export default class PathHelpService extends Service {
  getPaths(apiPath: string, backend: string, itemType?: string, itemID?: string): Promise<PathInfo>;
  getFormMetadata(helpUrl: string): Promise<FormMetadata>;
}
