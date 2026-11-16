/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { FormOptions } from 'vault/forms/form';

interface PathFilterConfigData {
  mode: string | null;
  paths: string[];
}

/** Keeps the secondary ID outside the payload and isolates editable filter paths. */
export default class PathFilterConfigForm extends Form<PathFilterConfigData> {
  declare readonly id?: string;
  declare mode: PathFilterConfigData['mode'];
  declare paths: PathFilterConfigData['paths'];

  formFields = [new FormField('mode', 'string'), new FormField('paths', 'array')];

  constructor(data: Partial<PathFilterConfigData> & { id?: string } = {}, options: FormOptions = {}) {
    const { id, mode = null, paths } = data;
    super({ mode, paths: Array.isArray(paths) ? [...paths] : [] }, options);
    this.id = id;
  }
}
