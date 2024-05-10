/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import FlashMessages from 'ember-cli-flash/services/flash-messages';

/*
we extend the ember-cli-flash service here so each ember engine can
import 'flash-messages' as a dependency giving it access to the
<FlashMessage> template in the main app's cluster.hbs file
*/
export default class FlashMessageService extends FlashMessages {}
