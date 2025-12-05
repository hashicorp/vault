/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

// Argument references for a TFVP (terraform vault provider) resource
interface TerraformOptions {
  name?: string;
  namespace?: string;
  policy?: string;
}
