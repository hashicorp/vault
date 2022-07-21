import { helper } from '@ember/component/helper';

export default helper(function ([a, ...toSubtract]) {
  return toSubtract.reduce((total, value) => total - parseInt(value, 0), a);
});
