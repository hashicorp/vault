import { helper as buildHelper } from '@ember/component/helper';
import { isNone, typeOf } from '@ember/utils';

export function reduceToArray(params) {
  return params.reduce(function(result, param) {
    if (isNone(param)) {
      return result;
    }
    if (typeOf(param) === 'array') {
      return result.concat(param);
    } else {
      return result.concat([param]);
    }
  }, []);
}

export default buildHelper(reduceToArray);
