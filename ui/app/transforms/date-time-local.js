import Transform from '@ember-data/serializer/transform';

export default class DateTimeLocalTransform extends Transform {
  getISODateFormat(deserializedDate) {
    // if the date is a date object or in local date time format ("yyyy-MM-dd'T'HH:mm"), we want to ensure
    // it gets converted to an ISOString
    if (
      typeof deserializedDate === 'object' ||
      (typeof deserializedDate === 'string' && !deserializedDate.includes('Z'))
    ) {
      return new Date(deserializedDate).toISOString();
    }

    return deserializedDate;
  }

  deserialize(serialized) {
    return serialized;
  }

  serialize(deserialized) {
    return this.getISODateFormat(deserialized);
  }
}
