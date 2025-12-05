/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export interface TerraformTemplateArgs {
  resource: string; // The Terraform resource type (e.g., "vault_auth_backend", "vault_mount")
  localId?: string; // The local identifier/label for this resource instance. Used in Terraform commands like `terraform import vault_auth_backend.example github` (<- is "example" here)
  options: TerraformOptions; //  Key/value pairs that build the terraform configuration
}

/**
 * Generates a Terraform resource block for Vault providers.
 * @see https://registry.terraform.io/providers/hashicorp/vault/latest/docs
 */
export const terraformTemplate = ({
  resource = '',
  localId = '<local identifier>',
  options,
}: TerraformTemplateArgs) => {
  const formattedContent = formatTerraformOptions(options);
  return `resource "${resource}" "${localId}" {
${formattedContent.join('\n\n')}
}`;
};

export const formatTerraformOptions = (options: TerraformOptions) => {
  const argReferences = [];
  for (const [key, value] of Object.entries(options)) {
    // Additional spaces before "key" are so argument references are indented over
    argReferences.push(`  ${key} = ${value}`);
  }
  return argReferences;
};
