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
    let uuid;
    // check if we can use crypto
    if (isSecureContext) {
      uuid = crypto.randomUUID();
      // odds are incredibly low that we'll run into a duplicate using crypto.randomUUID()
      // but just to be safe...
      while (this.ids.has(uuid)) {
        uuid = crypto.randomUUID();
      }
    } else {
      // fallback to incrementing id based on set size
      // initial size will be 0 so first id will be 1
      uuid = this.ids.size + 1;
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
