# Terraform Mapping Generator

The `generate:terraform-mapping` script contains TypeScript code that connects a Vault API feature to its Terraform provider resource. The output is a typed mapping function that `terraform-registry.ts` calls to render an HCL snippet in the UI's Automation Snippets component.

You only need to run this when adding Terraform snippet support to a new feature. If the feature already has a mapping in `app/utils/terraform-mappings/`, you do not need to run it again.

---

## Contents

- [Prerequisites](#prerequisites)
- [Two modes](#two-modes)
  - [OpenAPI mode - recommended](#openapi-mode---recommended)
  - [Interactive mode - manual fallback](#interactive-mode---manual-fallback)
- [Output](#output)
- [After generating](#after-generating)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

- `pnpm install` has been run (installs `tsx`, which the script needs to execute TypeScript)
- For **OpenAPI mode** only: `terraform` CLI installed (`brew install terraform`)

---

## Two modes

### OpenAPI mode - recommended

Add an entry to `app/utils/terraform-code-generator/terraform-resource-map.ts` mapping your API path to the Terraform resource type:

```ts
export const TERRAFORM_RESOURCE_MAP: Record<string, string> = {
  '/sys/policies/acl/{name}': 'vault_policy',
  '/sys/your/path/{name}': 'vault_your_resource', // ← add this
};
```

You can find the correct resource name in the [Terraform Vault provider docs](https://registry.terraform.io/providers/hashicorp/vault/latest/docs).

Pass the Vault API path as the argument. The script reads field definitions from the bundled Vault OpenAPI spec and the live Terraform provider schema, cross-references them, and emits a fully typed scaffold.

```
pnpm generate:terraform-mapping /sys/policies/acl/{name}
```

**What happens:**

1. Looks up the Terraform resource for your path in the resource map
2. Reads field definitions from `@hashicorp/vault-client-typescript/openapi.json`
3. Fetches the Terraform provider schema (downloads the provider on first run takes ~30 seconds)
4. Cross-references the two sources: fields that exist in both become the scaffold's payload interface; fields only in one source are noted in comments
5. Writes a `.ts` file to `app/utils/terraform-mappings/`

If any prerequisite is missing (no resource map entry, no OpenAPI spec, no `terraform` CLI), the script falls back to interactive mode automatically.

---

### Interactive mode - manual fallback

Pass a camelCase method name, or omit the argument entirely. The script prompts you for each piece of information.

```
pnpm generate:terraform-mapping mountsEnableSecretsEngine
```

You will be prompted for:

- **Terraform resource type** — e.g. `vault_mount`
- **Registry feature key** — the key used in `terraform-registry.ts`, e.g. `secrets/kv`
- **Fields** — name, type, and whether required, one at a time; press enter with no name when done

Supported field types: `string`, `boolean`, `number`, `object`, `array`, `heredoc`

Use `heredoc` for multi-line string fields like `policy` or `rules` — the generator will wrap them with `formatEot()` automatically.

---

## Output

Both modes write a file to `app/utils/terraform-mappings/` and run Prettier on it. Example output filename:

```
app/utils/terraform-mappings/sys-policies-acl-name-mapping.ts
```

The file contains:

- A TypeScript interface for the payload
- A mapping function that calls `terraformResourceTemplate`
- A commented-out registry snippet to paste into `terraform-registry.ts`
- `// TODO` comments on any `object` or `array` fields that need manual handling

---

## After generating

1. **Review the scaffold** — remove any fields not relevant to your feature
2. **Resolve `// TODO` items** — `object` and `array` fields need manual implementation
3. **Verify field names** against the [Terraform provider docs](https://registry.terraform.io/providers/hashicorp/vault/latest/docs)
4. **Register the mapping** — copy the commented registry snippet from the bottom of the generated file into `app/utils/terraform-registry.ts`

---

## Troubleshooting

**`ERR_UNKNOWN_FILE_EXTENSION .ts`** — `tsx` is not installed. Run `pnpm install`.

**`No Terraform resource found for "…"`** — add the path to `terraform-resource-map.ts` (see OpenAPI mode above).

**`terraform not found in PATH`** — install the Terraform CLI: `brew install terraform`. The script falls back to interactive mode if it's missing.

**`Path "…" not found in OpenAPI spec`** — the path may be enterprise-only or not yet in the bundled spec. Use interactive mode instead.
