export default {
  getItem(key) {
    var item = window.localStorage.getItem(key);
    return item && JSON.parse(item);
  },

  setItem(key, val) {
    window.localStorage.setItem(key, JSON.stringify(val));
  },

  removeItem(key) {
    return window.localStorage.removeItem(key);
  },

  keys() {
    return Object.keys(window.localStorage);
  },

  cleanUpStorage(string, keyToKeep) {
    if (!string) return;
    const relevantKeys = this.keys().filter((str) => str.startsWith(string));
    relevantKeys?.forEach((key) => {
      if (key !== keyToKeep) {
        localStorage.removeItem(key);
      }
    });
  },
};
