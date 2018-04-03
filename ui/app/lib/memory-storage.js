let cache = {};

export default {
  getItem(key) {
    var item = cache[key];
    return item && JSON.parse(item);
  },

  setItem(key, val) {
    cache[key] = JSON.stringify(val);
  },

  removeItem(key) {
    delete cache[key];
  },

  keys() {
    return Object.keys(cache);
  },
};
