import { helper } from '@ember/component/helper';
import { debug } from '@ember/debug';
import { htmlSafe } from '@ember/template';
import { sanitize } from 'dompurify';

export default helper(function sanitizedHtml([htmlString]) {
  try {
    return htmlSafe(sanitize(htmlString));
  } catch (e) {
    debug('Error sanitizing string', e);
    return '';
  }
});
