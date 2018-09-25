import { helper as buildHelper } from '@ember/component/helper';

export function jsonify([target]) {
  return JSON.parse(target);
}

export default buildHelper(jsonify);
