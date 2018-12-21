import DateBase from './-date-base';

export default DateBase.extend({
  compute() {
    this._super(...arguments);

    return Date.now();
  },
});
