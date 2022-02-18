export default {
  getItem(key) {
    var item = window.sessionStorage.getItem(key);
    return item && JSON.parse(item);
  },

  setItem(key, val) {
    window.sessionStorage.setItem(key, JSON.stringify(val));
  },

  removeItem(key) {
    return window.sessionStorage.removeItem(key);
  },

  keys() {
    return Object.keys(window.sessionStorage);
  },
};
