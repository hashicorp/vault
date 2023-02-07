/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// manage a set of unique ids
export default class {
  constructor() {
    this.ids = new Set();
  }
  /**
   * Returns a unique identifier.
   *
   * @method fetch
   * @param {Object} data Records attributes hash
   * @return {String} Unique identifier
   * @public
   */
  fetch() {
    let uuid = crypto.randomUUID();
    // odds are incredibly low that we'll run into a duplicate using crypto.randomUUID()
    // but just to be safe...
    while (this.ids.has(uuid)) {
      uuid = crypto.randomUUID();
    }
    this.ids.add(uuid);
    return uuid;
  }
  /**
   * Register an identifier.
   * Must throw if identifier is already used.
   *
   * @method set
   * @param {String|Number} id
   * @public
   */
  set(id) {
    if (this.ids.has(id)) {
      throw new Error(`ID ${id} is in use.`);
    }
    this.ids.add(id);
  }
  /**
   * Reset identity manager.
   *
   * @method reset
   * @public
   */
  reset() {
    this.ids.clear();
  }
}
