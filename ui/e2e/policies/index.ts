/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import fs from 'fs';
import path from 'path';

const readFile = (filePath: string) => {
  return fs.readFileSync(path.join(__dirname, filePath), 'utf-8');
};

export const USER_POLICY_MAP = {
  superuser: readFile('./superuser.hcl'),
};
