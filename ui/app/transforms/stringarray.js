import Transform from '@ember-data/serializer/transform';

export default class StringarrayTransform extends Transform {
  deserialize(serialized) {
    // client expects array of strings
    return serialized.split(',');
  }

  serialize(deserialized) {
    // api expects string with commas
    return deserialized.join(',');
  }
}
