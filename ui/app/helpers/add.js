import { helper as buildHelper } from '@ember/component/helper';

export function add(params) {
  let paramsFormatted = params.map(param => {
    return isNaN(param) ? 0 : param;
  });

  return paramsFormatted.reduce((sum, param) => parseInt(param, 0) + sum, 0);
}

export default buildHelper(add);
