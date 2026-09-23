/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';

export interface NamespaceFormData {
  path?: string;
}

export default class NamespaceForm extends Form<NamespaceFormData> {
  validations: Validations = {
    path: [
      { type: 'presence', message: `Path can't be blank.` },
      { type: 'endsInSlash', message: `Path can't end in forward slash '/'.` },
      { type: 'containsWhiteSpace', message: `Path can't contain whitespace.` },
    ],
  };

  formFields = [new FormField('path', 'string')];
}
