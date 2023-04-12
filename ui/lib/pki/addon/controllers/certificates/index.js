/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { getOwner } from '@ember/application';
import { action } from '@ember/object';

export default class PkiCertificatesIndexController extends Controller {
  queryParams = ['pageFilter'];

  get filteredCertificates() {
    const pageFilter = this.pageFilter;
    const certificates = this.model.certificates;

    if (pageFilter) {
      return certificates.filter((cert) => cert.id.toLowerCase().includes(pageFilter.toLowerCase()));
    } else {
      return certificates;
    }
  }

  get mountPoint() {
    return getOwner(this).mountPoint;
  }

  @action setFilter(val) {
    this.filter = val;
  }
  @action setFilterFocus(bool) {
    this.filterFocused = bool;
  }
}
