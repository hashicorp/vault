import ApplicationAdapter from './application';

const WRAPPING_ENDPOINTS = ['lookup', 'wrap', 'unwrap', 'rewrap'];
const TOOLS_ENDPOINTS = ['random', 'hash'];

export default ApplicationAdapter.extend({
  toolUrlFor(action) {
    const isWrapping = WRAPPING_ENDPOINTS.includes(action);
    const isTool = TOOLS_ENDPOINTS.includes(action);
    const prefix = isWrapping ? 'wrapping' : 'tools';
    if (!isWrapping && !isTool) {
      throw new Error(`Calls to a ${action} endpoint are not currently allowed in the tool adapter`);
    }
    return `${this.buildURL()}/${prefix}/${action}`;
  },

  toolAction(action, data, options = {}) {
    const { wrapTTL } = options;
    const url = this.toolUrlFor(action);
    const ajaxOptions = wrapTTL ? { data, wrapTTL } : { data };
    return this.ajax(url, 'POST', ajaxOptions);
  },
});
