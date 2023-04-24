import Transform from '@ember-data/serializer/transform';

export default class TtlTransform extends Transform {
  deserialize(serialized, { enabledKey, durationKey, isOppositeValue }) {
    return {
      enabled: isOppositeValue ? !serialized[enabledKey] : serialized[enabledKey],
      duration: serialized[durationKey],
    };
  }

  serialize({ enabled, duration }, { enabledKey, durationKey, isOppositeValue }) {
    return {
      [enabledKey]: isOppositeValue ? !enabled : enabled,
      [durationKey]: duration,
    };
  }
}
