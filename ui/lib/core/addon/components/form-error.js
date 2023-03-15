/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * @module FormError
 * FormError components are used to show an error on a form field that is more compact than the
 * normal MessageError component. This component adds an icon and styling to the content of the
 * component, so additionally styling (bold, italic) and links are allowed.
 *
 * @example
 * ```js
 * <FormError>Oh no <em>something bad</em>! <a href="#">Do something</a></FormError>
 * ```
 */

import Component from '@ember/component';
import layout from '../templates/components/form-error';

export default Component.extend({
  layout,
});
