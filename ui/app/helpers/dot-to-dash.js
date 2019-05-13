import { helper as buildHelper } from '@ember/component/helper';

export function dotToDash([string]) {
  return string.replace(/\./gi, '-');
}

export default buildHelper(dotToDash);
