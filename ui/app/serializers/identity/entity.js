import DS from 'ember-data';
import IdentitySerializer from './_base';

export default IdentitySerializer.extend(DS.EmbeddedRecordsMixin, {
  // we don't need to serialize relationships here
  serializeHasMany() {},
  attrs: {
    aliases: { embedded: 'always' },
  },
});
