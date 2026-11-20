/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form, { type FormOptions } from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import { keyIsFolder, parentKeyForKey } from 'core/utils/key-utils';
import type { Validations } from 'vault/app-types';

type SecretFormData = Record<string, unknown> & {
  backend?: string;
  canDelete?: boolean;
  failedServerRead?: boolean;
  id?: string;
  path?: string;
  secretData?: Record<string, unknown>;
};

/** Holds KV v1 secret form state and serializes the request body. */
export default class SecretForm extends Form<SecretFormData> {
  declare path: string | undefined;
  pathAttr = 'path';
  formFields = [new FormField('path', 'string'), new FormField('secretData', 'object')];
  validations: Validations = {};

  constructor(data: Partial<SecretFormData> = {}, options: FormOptions = {}) {
    super(data, options);
    this.validations = options.isNew ? { path: [{ type: 'presence', message: "Path can't be blank." }] } : {};
  }

  get id() {
    return this.data.id ?? this.data.path ?? '';
  }

  get backend() {
    return this.data.backend;
  }

  get canDelete() {
    return this.data.canDelete;
  }

  get failedServerRead() {
    return this.data.failedServerRead;
  }

  get secretKeyAndValue() {
    return Object.keys(this.data.secretData || {}).map((key) => ({
      key,
      value: this.data.secretData?.[key],
    }));
  }

  get isFolder() {
    return keyIsFolder(this.path || this.id);
  }

  get parentKey() {
    return parentKeyForKey(this.path || this.id);
  }

  toJSON() {
    const form = super.toJSON();
    return {
      ...form,
      data: { ...(this.data.secretData || {}) },
    };
  }
}
