/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  sanitizeId,
  terraformResourceTemplate,
  terraformVariableTemplate,
} from 'core/utils/code-generators/terraform';

export const generateTerraformSnippet = (namespacePaths: string[], currentPath: string): string => {
  const tfTopResources: string[] = [];
  const tfMiddleVariables: string[] = [];
  const tfMiddleResources: string[] = [];
  const tfBottomVariables: string[] = [];
  const tfBottomResources: string[] = [];

  // Parse  to group by hierarchy
  const topLevels = new Set<string>();
  const middleLayers: { [topLayer: string]: Set<string> } = {};
  const bottomLayers: { [middleKey: string]: Set<string> } = {};

  namespacePaths?.forEach((nsPath) => {
    const parts = nsPath.split('/');
    const topLayer = parts[0] as string;

    if (parts.length === 1) {
      topLevels.add(topLayer);
    } else if (parts.length === 2) {
      topLevels.add(topLayer);
      if (!middleLayers[topLayer]) middleLayers[topLayer] = new Set();
      middleLayers[topLayer].add(parts[1] as string);
    } else if (parts.length === 3) {
      topLevels.add(topLayer);
      const middleLayer = parts[1] as string;
      const bottomLayer = parts[2] as string;
      if (!middleLayers[topLayer]) middleLayers[topLayer] = new Set();
      middleLayers[topLayer].add(middleLayer);
      const middleKey = `${topLayer}/${middleLayer}`;
      if (!bottomLayers[middleKey]) bottomLayers[middleKey] = new Set();
      bottomLayers[middleKey].add(bottomLayer);
    }
  });

  // Generate Terraform resources
  topLevels.forEach((topLayer) => {
    const sanitizedTopId = sanitizeId(topLayer);

    // Top layer resource
    const topResourceArgs: { [key: string]: string } = { path: `"${topLayer}"` };
    if (currentPath) {
      topResourceArgs['namespace'] = `"${currentPath}"`;
    }

    tfTopResources.push(
      terraformResourceTemplate({
        resource: 'vault_namespace',
        localId: sanitizedTopId,
        resourceArgs: topResourceArgs,
      })
    );

    // Middle layers for this top layer
    const middles = middleLayers[topLayer];
    if (middles && middles.size > 0) {
      const middleChildren = Array.from(middles).map((m) => `"${m}"`);

      // Middle variable and resource
      tfMiddleVariables.push(
        terraformVariableTemplate({
          variable: `${sanitizedTopId}_child_namespaces`,
          variableArgs: {
            type: 'set(string)',
            default: `[${middleChildren.join(', ')}]`,
          },
        })
      );

      const namespaceReference = currentPath
        ? `vault_namespace.${sanitizedTopId}.path_fq`
        : `vault_namespace.${sanitizedTopId}.path`;

      tfMiddleResources.push(
        terraformResourceTemplate({
          resource: 'vault_namespace',
          localId: `${sanitizedTopId}_children`,
          resourceArgs: {
            for_each: `var.${sanitizedTopId}_child_namespaces`,
            namespace: namespaceReference,
            path: 'each.key',
          },
        })
      );

      // Bottom layers for each middle layer
      middles.forEach((middleLayer) => {
        const middleKey = `${topLayer}/${middleLayer}`;
        const bottoms = bottomLayers[middleKey];

        if (bottoms && bottoms.size > 0) {
          const sanitizedMiddleId = sanitizeId(middleLayer);
          const bottomChildren = Array.from(bottoms).map((b) => `"${b}"`);

          // Bottom variable and resource
          tfBottomVariables.push(
            terraformVariableTemplate({
              variable: `${sanitizedTopId}_${sanitizedMiddleId}_child_namespaces`,
              variableArgs: {
                type: 'set(string)',
                default: `[${bottomChildren.join(', ')}]`,
              },
            })
          );

          tfBottomResources.push(
            terraformResourceTemplate({
              resource: 'vault_namespace',
              localId: `${sanitizedTopId}_${sanitizedMiddleId}_children`,
              resourceArgs: {
                for_each: `var.${sanitizedTopId}_${sanitizedMiddleId}_child_namespaces`,
                namespace: `vault_namespace.${sanitizedTopId}_children["${middleLayer}"].path_fq`,
                path: 'each.key',
              },
            })
          );
        }
      });
    }
  });

  // Build in proper dependency order
  const orderedSections = [
    tfMiddleVariables.join('\n\n'),
    tfBottomVariables.join('\n\n'),
    tfTopResources.join('\n\n'),
    tfMiddleResources.join('\n\n'),
    tfBottomResources.join('\n\n'),
  ].filter((section) => section.trim() !== '');

  return orderedSections.join('\n\n');
};

export const generateCliSnippet = (namespacePaths: string[], currentPath: string): string => {
  const cliSnippet: string[] = [];

  namespacePaths.forEach((nsPath) => {
    const parts = nsPath.split('/');

    if (parts.length === 1) {
      // Top level namespace
      const fullPath = currentPath ? `-namespace ${currentPath} ` : '';
      cliSnippet.push(`vault namespace create ${fullPath}${nsPath}/`);
    } else if (parts.length === 2) {
      // Middle level namespace
      const parentNs = parts[0];
      const fullPath = currentPath ? `${currentPath}/${parentNs}` : parentNs;
      cliSnippet.push(`vault namespace create -namespace ${fullPath} ${parts[1]}/`);
    } else if (parts.length === 3) {
      // Bottom level namespace
      const parentNs = parts[0] + '/' + parts[1];
      const fullPath = currentPath ? `${currentPath}/${parentNs}` : parentNs;
      cliSnippet.push(`vault namespace create -namespace ${fullPath} ${parts[2]}/`);
    }
  });

  return cliSnippet.join('\n');
};

export const generateApiSnippet = (namespacePaths: string[], currentPath: string): string => {
  const apiSnippet: string[] = namespacePaths.map((nsPath) => {
    const parts = nsPath.split('/');

    if (parts.length === 1) {
      const nsHeader = currentPath ? `    --header "X-Vault-Namespace: ${currentPath}"\\\n` : '';
      return `curl \\\n    --header "X-Vault-Token: $VAULT_TOKEN" \\\n${nsHeader}    --request PUT \\\n    $VAULT_ADDR/v1/sys/namespaces/${nsPath}`;
    } else {
      const parentPath = currentPath + '/' + parts.slice(0, -1).join('/');
      const childName = parts[parts.length - 1];
      return `curl \\\n    --header "X-Vault-Token: $VAULT_TOKEN" \\\n    --header "X-Vault-Namespace: ${parentPath}" \\\n    --request PUT \\\n    $VAULT_ADDR/v1/sys/namespaces/${childName}`;
    }
  });

  return apiSnippet.join('\n\n');
};
