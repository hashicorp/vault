/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationSerializer from '../application';

const MAP_ATTRS_TO_MODEL = {
  crl_expiry_data: {
    enabledKey: 'disable',
    isOppositeValue: true,
    durationKey: 'expiry',
  },
  auto_rebuild_data: {
    enabledKey: 'auto_rebuild',
    isOppositeValue: false,
    durationKey: 'auto_rebuild_grace_period',
  },
  delta_crl_building_data: {
    enabledKey: 'enable_delta',
    isOppositeValue: false,
    durationKey: 'delta_rebuild_interval',
  },
  ocsp_expiry_data: {
    enabledKey: 'ocsp_disable',
    isOppositeValue: true,
    durationKey: 'ocsp_expiry',
  },
};
export default class PkiCrlSerializer extends ApplicationSerializer {
  normalize(model, data) {
    for (const key in MAP_ATTRS_TO_MODEL) {
      const { enabledKey, durationKey, isOppositeValue } = MAP_ATTRS_TO_MODEL[key];
      const valueBlock = {
        enabled: isOppositeValue ? !data[enabledKey] : data[enabledKey],
        duration: data[durationKey],
      };
      data = { ...data, [key]: valueBlock };
      delete data[enabledKey];
      delete data[durationKey];
    }
    return super.normalize(...arguments);
  }

  serialize() {
    const json = super.serialize(...arguments);
    for (const key in MAP_ATTRS_TO_MODEL) {
      if (key in json) {
        const { enabledKey, durationKey, isOppositeValue } = MAP_ATTRS_TO_MODEL[key];
        const { enabled, duration } = json[key];
        json[enabledKey] = isOppositeValue ? !enabled : enabled;
        json[durationKey] = duration;
        delete json[key];
      }
    }
    return json;
  }
}
