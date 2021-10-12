import { helper } from '@ember/component/helper';

export default helper(function isEmptyValue([value] /*, hash*/) {
  if (typeof value === 'object' && value !== null) {
    return Object.keys(value).length === 0;
  }
  return value == null || value === '';
});
