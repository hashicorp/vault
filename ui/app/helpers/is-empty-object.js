import { helper } from '@ember/component/helper';

export default helper(function isEmptyObject([object] /*, hash*/) {
  return Object.keys(object).length === 0;
});
