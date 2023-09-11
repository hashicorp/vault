import Transform from '@ember-data/serializer/transform';

/**
 * correctly maps boolean values to the two options for listingVisibility
 * attribute on seceret engines and auth engines
 */
export default class MountVisibilityTransform extends Transform {
  deserialize(serialized) {
    return serialized === 'unauth';
  }

  serialize(deserialized) {
    return deserialized === true ? 'unauth' : 'hidden';
  }
}
