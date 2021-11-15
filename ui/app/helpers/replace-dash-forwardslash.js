import { helper as buildHelper } from '@ember/component/helper';

export function replaceDashForwardslash([string]) {
  if (!string) {
    return;
  }
  string = string.replace(/-/g, '/');
  return string;
}

export default buildHelper(replaceDashForwardslash);
