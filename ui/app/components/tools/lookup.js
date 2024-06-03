/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { addSeconds, parseISO } from 'date-fns';

/**
 * @module ToolLookup
 * ToolLookup components are components that sys/wrapping/lookup functionality.  Most of the functionality is passed through as actions from the tool-actions-form and then called back with properties.
 *
 * @example
 * <Tools::Lookup @creation_time={{creation_time}} @creation_ttl={{creation_ttl}} @creation_path={{creation_path}} @token={{token}} @onClear={{action "onClear"}} @errors={{errors}}/>
 *
 * @param {string} creation_time - ISO string creation time for token
 * @param {number} creation_ttl - token ttl in seconds
 * @param {string} creation_path - path where token was originally generated
 * @param {string} token=null - token
 * @param {function} onClear - callback that resets all of values to defaults. Must be passed as `{{action "onClear"}}`
 * @param {object} errors=null - errors returned if lookup request fails
 */
export default class ToolLookup extends Component {
  get expirationDate() {
    // returns new Date with seconds added.
    return addSeconds(parseISO(this.args.creation_time), this.args.creation_ttl);
  }
}
