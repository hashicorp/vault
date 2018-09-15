/*jshint eqeqeq: false */
import { helper as buildHelper } from '@ember/component/helper';

export function coerceEq(params) {
  return params[0] == params[1];
}

export default buildHelper(coerceEq);
