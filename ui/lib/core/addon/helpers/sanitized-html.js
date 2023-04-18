import { helper } from '@ember/component/helper';
import { debug } from '@ember/debug';
import { htmlSafe } from '@ember/template';
import { sanitize } from 'dompurify';

export default helper(function sanitizedHtml([htmlString]) {
  try {
    return htmlSafe(sanitize(htmlString));
  } catch (e) {
    debug('Error sanitizing string', e);
    // I couldn't get this to actually fail but as a fallback,
    // render the value as-is with HTML escaping
    return htmlString;
  }
});
