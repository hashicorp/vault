import { helper } from '@ember/component/helper';
import { WIF_ENGINES } from 'vault/helpers/mountable-secret-engines';

// a helper to use within templates to determine if a secret engine is a WIF engine
// we cannot use the mountable-secret-engines helper for this purpose because the exported helper on that file is for mountableEngines. WIF methods cannot be accessed from within template.
export default helper(function isWifEngine([type]) {
  return WIF_ENGINES.includes(type);
});
