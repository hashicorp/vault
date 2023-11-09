/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import templateOnly from '@ember/component/template-only';

/**
 * @module FormFieldLabel
 * FormFieldLabel components add labels and descriptions to inputs
 *
 * @example
 * <FormFieldLabel for="input-name" @label={{this.label}} @helpText={{this.helpText}} @subText={{this.subText}} @docLink={{this.docLink}} />
 *
 * @param {string} [label] - label text -- component attributes are spread on label element
 * @param {string} [helpText] - adds a tooltip
 * @param {string} [subText] - de-emphasized text rendered below the label
 * @param {string} [docLink] - url to documentation rendered after the subText
 */

export default templateOnly();
