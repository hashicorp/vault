import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    payload.data.name = payload.id;
    if (payload.data.alphabet) {
      payload.data.alphabet = [payload.data.alphabet];
    }
    // strip out P character from any named capture groups
    if (payload.data.pattern) {
      this._formatNamedCaptureGroups(payload.data, '?P', '?');
    }
    return this._super(store, primaryModelClass, payload, id, requestType);
  },

  serialize() {
    const json = this._super(...arguments);
    if (json.alphabet && Array.isArray(json.alphabet)) {
      // Templates should only ever have one alphabet
      json.alphabet = json.alphabet[0];
    }
    // add P character to any named capture groups
    if (json.pattern) {
      this._formatNamedCaptureGroups(json, '?', '?P');
    }
    return json;
  },

  _formatNamedCaptureGroups(json, replace, replaceWith) {
    // named capture groups are handled differently between Go and js
    // first look for named capture groups in pattern string
    const regex = new RegExp(/\?P?(<(.+?)>)/, 'g');
    const namedGroups = json.pattern.match(regex);
    if (namedGroups) {
      namedGroups.forEach((group) => {
        // add or remove P depending on destination
        json.pattern = json.pattern.replace(group, group.replace(replace, replaceWith));
      });
    }
  },

  extractLazyPaginatedData(payload) {
    return payload.data.keys.map((key) => {
      const model = {
        id: key,
        name: key,
      };
      if (payload.backend) {
        model.backend = payload.backend;
      }
      return model;
    });
  },
});
