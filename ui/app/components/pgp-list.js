import Component from '@ember/component';

export default Component.extend({
  onDataUpdate: () => {},
  listLength: 0,
  listData: null,

  init() {
    this._super(...arguments);
    let num = this.listLength;
    if (num) {
      num = parseInt(num, 10);
    }
    const list = this.newList(num);
    this.set('listData', list);
  },

  didReceiveAttrs() {
    this._super(...arguments);
    let list;
    if (!this.listLength) {
      this.set('listData', []);
      return;
    }
    // no update needed
    if (this.listData.length === this.listLength) {
      return;
    }
    if (this.listLength < this.listData.length) {
      // shorten the current list
      list = this.listData.slice(0, this.listLength);
    } else if (this.listLength > this.listData.length) {
      // add to the current list by creating a new list and copying over existing list
      list = [...this.listData, ...this.newList(this.listLength - this.listData.length)];
    }
    this.set('listData', list || this.listData);
    this.onDataUpdate((list || this.listData).compact().map((k) => k.value));
  },

  newList(length) {
    return Array(length || 0)
      .fill(null)
      .map(() => ({ value: '' }));
  },

  actions: {
    setKey(index, key) {
      const { listData } = this;
      listData.splice(index, 1, key);
      this.onDataUpdate(listData.compact().map((k) => k.value));
    },
  },
});
