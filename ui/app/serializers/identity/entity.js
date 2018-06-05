import DS from 'ember-data';
import IdentitySerializer from './_base';

export default IdentitySerializer.extend(DS.EmbeddedRecordsMixin, {
  attrs: {
    aliases: { embedded: 'always' },
  },
});
