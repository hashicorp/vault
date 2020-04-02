!(function define(Object) {
  Object
    ? typeof globalThis == 'object' ||
      Object.prototype.__defineGetter__('_', define) ||
      // eslint-disable-next-line no-undef
      _ ||
      delete Object.prototype._
    : (this.globalThis = this)
})(Object)
