/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { sanitizePath, sanitizeStart } from 'core/utils/sanitize-path';
import { task } from 'ember-concurrency';

/**
 * PermissionsService
 * ---------------------------------------------------------------------------
 * What it does
 * - Gates sidebar visibility and the “Resultant ACL” banner.
 * - Consumes one resultant-acl payload: { exact_paths, glob_paths, root, chroot_namespace }.
 * - Evaluates permissions against the fully-qualified path:
 *   <chroot>/<currentNamespace>/<apiPath> (omitting empty parts).
 *
 * Policy semantics (brief)
 * - '+' = exactly one segment; '*' = wildcard; '/' = include base; '/*' = children-only.
 * - Globpaths return policy paths with * if the path includes a wildcard. Meaning,
 *  sys/admin/* returns as sys/admin/, but +/sys/admin/* returns as +/sys/admin/*
 *  the implication is that a plus sign make the * a non-greedy match (e.g. children only), but without a plus
 *  a policy path is greedy (e.g. base + children)
 * - Most-specific (longest) match wins; any explicit 'deny' beats allow.
 *   Docs: https://developer.hashicorp.com/vault/docs/concepts/policies#priority-matching
 *
 * Sidebar flow
 * - Templates → has-permissions → hasNavPermission() → hasPermission() on resolved API paths.
 *
 * Banner flow
 * - permissionsBanner getter computes effective namespace and returns:
 *   • null if root, read-failed if ACL fetch failed, or if any canary path is allowed here
 *   • no-ns-access otherwise
 * - “Canary paths” are a small set of UI-centric endpoints we probe in the current namespace.
 * - They build off API_PATHS, which also drives sidebar nav items.
 * - If you ever wanted to add to a canary probe list, either expand API_PATHS
 * or add to CANARY_PATHS directly.
 *
 * Chroot note
 * - When chroot_namespace is present, backend keys are already prefixed.
 *   We mirror this by composing fullCurrentNamespace for all checks.
 *
 * Example resultant-acl payload:
 * {
 *   "exact_paths": {
 *     "sys/policies/acl": { "capabilities": ["read", "list"] },
 *     "sys/policies/acl/my-policy": { "capabilities": ["read"] },
 *     "sys/auth": { "capabilities": ["deny"] }
 *   },
 *   "glob_paths": {
 *     "secret/data/finance/+/payroll": { "capabilities": ["create", "update", "read", "list"] },
 *     "secret/data/engineering/*": { "capabilities": ["read", "list"] },
 *     "secret/data/hr/": { "capabilities": ["read", "list"] },
 *     "+/auth/*": { "capabilities": ["deny"] },
 *     "": { "capabilities": ["read", "list"] } // baseline allow-all
 *   },
 *   "root": false,
 *   "chroot_namespace": "ns1/child"
 * }
 */

export const PERMISSIONS_BANNER_STATES = {
  readFailed: 'read-failed',
  noAccess: 'no-ns-access',
};

export const RESULTANT_ACL_PATH = 'sys/internal/ui/resultant-acl'; // export for tests

const API_PATHS = {
  access: {
    methods: 'sys/auth',
    mfa: 'identity/mfa/method',
    oidc: 'identity/oidc/client',
    entities: 'identity/entity/id',
    groups: 'identity/group/id',
    leases: 'sys/leases/lookup',
    namespaces: 'sys/namespaces',
    'control-groups': 'sys/control-group/',
  },
  policies: {
    acl: 'sys/policies/acl',
    rgp: 'sys/policies/rgp',
    egp: 'sys/policies/egp',
  },
  tools: {
    wrap: 'sys/wrapping/wrap',
    lookup: 'sys/wrapping/lookup',
    unwrap: 'sys/wrapping/unwrap',
    rewrap: 'sys/wrapping/rewrap',
    random: 'sys/tools/random',
    hash: 'sys/tools/hash',
  },
  status: {
    replication: 'sys/replication',
    license: 'sys/license',
    seal: 'sys/seal',
    raft: 'sys/storage/raft/configuration',
  },
  clients: {
    activity: 'sys/internal/counters/activity',
    config: 'sys/internal/counters/config',
  },
  settings: {
    customMessages: 'sys/config/ui/custom-messages',
  },
  sync: {
    destinations: 'sys/sync/destinations',
    associations: 'sys/sync/associations',
    config: 'sys/sync/config',
    github: 'sys/sync/github-apps',
  },
  monitoring: {
    'utilization-report': 'sys/utilization-report',
  },
};

const API_PATHS_TO_ROUTE_PARAMS = {
  'sys/auth': { route: 'vault.cluster.access.methods', models: [] },
  'identity/entity/id': { route: 'vault.cluster.access.identity', models: ['entities'] },
  'identity/group/id': { route: 'vault.cluster.access.identity', models: ['groups'] },
  'sys/leases/lookup': { route: 'vault.cluster.access.leases', models: [] },
  'sys/namespaces': { route: 'vault.cluster.access.namespaces', models: [] },
  'sys/control-group/': { route: 'vault.cluster.access.control-groups', models: [] },
  'identity/mfa/method': { route: 'vault.cluster.access.mfa', models: [] },
  'identity/oidc/client': { route: 'vault.cluster.access.oidc', models: [] },
};

// Canary endpoints: quick check for “meaningful UI access” in the *current* namespace.
// If the token has any non-deny capability on any canary here, we suppress the banner.
// This does not try to cover all possible paths (e.g. secrets engines only such as `+/kv/data/*`),
// by design—probing everything is infeasible. Keep the list small and UI-centric.
//
// IMPORTANT NOTE: A user scoped only to say a secrets engine path
// (e.g. `+/kv/data/*`) would still trigger the banner, since iterating over all
// possible paths is not feasible.
//
// This may be the situation for Namespace-tenancy setups where users are
// confined to secrets engines in a child namespace and given no management of
// endpoints there.
const CANARY_PATHS = [
  ...Object.values(API_PATHS.access),
  ...Object.values(API_PATHS.policies),
  ...Object.values(API_PATHS.tools),
  ...Object.values(API_PATHS.status),
  ...Object.values(API_PATHS.clients),
  ...Object.values(API_PATHS.settings),
  ...Object.values(API_PATHS.sync),
  ...Object.values(API_PATHS.monitoring),
];

export default class PermissionsService extends Service {
  @tracked exactPaths = null;
  @tracked globPaths = null;
  @tracked hasFallbackAccess = false;
  @tracked isRoot = false;
  @tracked _aclLoadFailed = false;
  @tracked chrootNamespace = null;

  @service store;
  @service namespace;

  // isAclLoaded:
  // - True if we know the caller is actual root (resp.data.root === true),
  //   or if we’ve received any exact/glob paths from the resultant-acl payload.
  // - Starts false until the first ACL fetch.
  // - On fetch failure, we deliberately keep isRoot=false (since the backend
  //   didn’t confirm) and set hasFallbackAccess=true instead. This way the
  //   sidebar remains usable but the banner shows "read-failed".
  get isAclLoaded() {
    return this.isRoot || this.exactPaths !== null || this.globPaths !== null;
  }

  get fullCurrentNamespace() {
    const currentNs = this.namespace.path;
    return this.chrootNamespace
      ? `${sanitizePath(this.chrootNamespace)}/${sanitizePath(currentNs)}`
      : sanitizePath(currentNs);
  }

  // Banner logic: only show when we *know* the user has no meaningful UI access
  // in the current namespace. Precedence:
  // 1) Root token → never show a banner.
  // 2) ACL fetch failed (_aclLoadFailed) → show "read-failed" (sidebar may still
  //    render via hasFallbackAccess; the banner explains missing ACL data).
  // 3) ACLs not loaded yet (!isAclLoaded) → suppress banner (avoid flicker while loading).
  // 4) Explicit UI check: if the caller can READ "sys/internal/ui/resultant-acl" in the
  //    effective namespace, treat that as minimal UI access → suppress banner.
  //    History: This was the original intent (CE PR #23503). We temporarily removed reliance
  //    on this signal (CE PR #25256) because our old banner logic used string prefix checks
  //    and the resultant-acl payload did not preserve namespace-segment wildcards (`+`), so
  //    policies like `+/sys/...` or `+/+/sys/...` in child namespaces could be missed.
  //    Today the backend returns `+` in `glob_paths` and this service matches against the
  //    fully-qualified path using proper glob semantics, so this explicit check is reliable again.
  // 5) Heuristic canaries: if any curated UI-centric path is non-deny, suppress the banner.
  // 6) Otherwise → "no-ns-access".
  get permissionsBanner() {
    // 1) Real root never sees a banner
    if (this.isRoot) return null;

    // 2) Fetch failed → explain
    if (this._aclLoadFailed) return PERMISSIONS_BANNER_STATES.readFailed;

    // 3) Still loading → stay quiet
    if (!this.isAclLoaded) return null;

    // 4) Explicit resultant-acl read → minimal UI access confirmed
    if (this.hasPermission(RESULTANT_ACL_PATH, ['read'])) {
      return null;
    }

    // 5) Heuristic: any canary with any non-deny capability
    const anyAllowed = CANARY_PATHS.some((p) => this.hasPermission(p));

    // 6) Final decision
    return anyAllowed ? null : PERMISSIONS_BANNER_STATES.noAccess;
  }

  // Load resultant-ACL used by sidebar + banner.
  // Success:
  //   • Clear _aclLoadFailed, then hydrate via setPaths(resp).
  // Failure (network/403/etc):
  //   • Mark _aclLoadFailed so the banner explains missing ACL data.
  //   • Enable hasFallbackAccess so the sidebar remains usable (legacy behavior).
  //   • Clear exact/glob to avoid using stale ACL from a prior success.
  @task *getPaths() {
    try {
      const resp = yield this.store.adapterFor('permissions').query();
      this._aclLoadFailed = false;
      this.setPaths(resp);
    } catch (err) {
      this._aclLoadFailed = true;
      this.isRoot = false;
      this.hasFallbackAccess = true; // fallback so nav stays visible

      // Avoid stale ACL from a previous successful load
      this.exactPaths = null;
      this.globPaths = null;
    }
  }

  // Populate tracked state from a successful resultant-ACL response.
  //   • exact_paths / glob_paths → used for all permission checks
  //   • root (boolean) → authoritative “real root” from backend
  //   • hasFallbackAccess → always false on success (only true if ACL fetch fails)
  //   • chroot_namespace → prefix applied to all checks
  setPaths(resp) {
    this.exactPaths = resp.data.exact_paths;
    this.globPaths = resp.data.glob_paths;
    this.isRoot = resp.data.root; // true root per backend
    this.hasFallbackAccess = false; // nav gating shortcut (true root behaves like “show all”). Only set true if ACL fail.
    this.chrootNamespace = resp.data.chroot_namespace;
  }

  reset() {
    this.exactPaths = null;
    this.globPaths = null;
    this.hasFallbackAccess = false;
    this.chrootNamespace = null;
    this._aclLoadFailed = false;
  }

  // ===== Matching helpers (STRICT) ==========================================

  // _globKeyToRegex
  // Converts a Vault glob-style policy key into a regex:
  //   • Trailing "/"   → include base + any children
  //   • Trailing "/*"  → children-only (≥1 child; base does NOT match)
  //   • "+"            → exactly one segment
  //   • "*"            → greedy segment matcher
  //   • "" (empty key) → baseline allow-all, equivalent to `path "*" { … }`
  //
  // PERF NOTE: Consider caching compiled regexes if ACLs are large or called frequently.
  _globKeyToRegex(globKey) {
    // Empty key → allow-all
    if (globKey === '') {
      return /^.*$/;
    }

    const endsWithChildrenOnly = /\/\*$/.test(globKey); // "/*" → children-only
    const endsWithIncludeBase = !endsWithChildrenOnly && /\/$/.test(globKey); // "/" → include base
    const trimmed = globKey
      .replace(/\/\*$/, '') // strip "/*"
      .replace(/\/$/, ''); // strip trailing "/"

    const parts = trimmed.split('/').map((seg) => {
      if (seg === '+') return '[^/]+'; // exactly one segment
      if (seg === '*') return '.*'; // greedy tail matcher
      return seg.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'); // escape literal
    });

    if (endsWithChildrenOnly) {
      // require ≥1 child segment
      return new RegExp('^' + parts.join('/') + '/.+' + '$');
    }

    if (endsWithIncludeBase) {
      // match base OR any children
      return new RegExp('^' + parts.join('/') + '(?:/.*)?' + '$');
    }

    // exact match only
    return new RegExp('^' + parts.join('/') + '$');
  }

  // _matchExact
  // Example: needle = "sys/auth"
  //   • "sys/auth"        → matches directly
  //   • "sys/auth/methods" → chosen over "sys/auth" (more specific child)
  //
  // Resolves an "exact_paths" match with Vault policy semantics:
  //   • Normalize trailing "/" for both needle and keys
  //   • Exact equality takes precedence
  //   • Otherwise, choose the longest child key (most specific) under the base
  _matchExact(fullPath) {
    const exact = this.exactPaths;
    if (!exact) return null;

    const needle = fullPath.replace(/\/$/, '');

    // Build [orig, norm] pairs so we can match on norm but return orig
    const tuples = Object.keys(exact).map((k) => [k, k.replace(/\/$/, '')]);

    // Exact equality wins outright
    const eq = tuples.find(([, norm]) => norm === needle);
    if (eq) {
      const [orig] = eq;
      return { key: orig, entry: exact[orig] };
    }

    // Otherwise, pick the longest child (most specific) under the base
    const child = tuples
      .filter(([, norm]) => norm.startsWith(needle + '/'))
      .sort((a, b) => b[1].length - a[1].length)[0];

    return child ? { key: child[0], entry: exact[child[0]] } : null;
  }

  // _matchGlob
  // Example: fullPath = "sys/auth/methods"
  //   • glob "*"                → matches (baseline allow-all)
  //   • glob "+/sys/*"            → matches, less specific
  //   • glob "+/sys/auth/*"       → matches, more specific → chosen
  //   • glob "sys/auth/methods" → exact equality handled in _matchExact
  //
  // Glob match precedence:
  //   • Empty key "" is a special baseline = allow all (overridable by denies)
  //   • Among regex matches, the longest key (most specific) wins
  //   • Final allow/deny resolution is deferred to _decide()
  _matchGlob(fullPath) {
    const globPaths = this.globPaths;
    if (!globPaths) return null;

    let matchKey = null;

    // Treat empty-root ('') as baseline "allow all" if present
    if (Object.prototype.hasOwnProperty.call(globPaths, '')) {
      matchKey = '';
    }

    for (const k of Object.keys(globPaths)) {
      if (k === '') continue; // already accounted for
      const re = this._globKeyToRegex(k);
      if (re.test(fullPath)) {
        if (matchKey === null || k.length > matchKey.length) {
          matchKey = k; // longest wins
        }
      }
    }

    return matchKey !== null ? { key: matchKey, entry: globPaths[matchKey] } : null;
  }

  // _decide
  // Resolution rule across matchers:
  // - If any matched entry is 'deny' → false (authoritative).
  // - Otherwise, ALL requested capabilities must be satisfied by at least one matched non-deny entry.
  //   • The 'capabilities' arg comes from the caller (via hasPermission):
  //     - ["list"] for identity pages
  //     - ["read", "update"] if a UI flow requires those
  //     - [null] means "any non-deny capability is enough" (common for sidebar/nav checks).
  _decide(fullPath, capabilities = [null]) {
    // Collect matches from BOTH exact and glob, then apply precedence:
    // 1. If ANY matching entry is an explicit deny → false.
    // 2. If no deny, allow if every requested capability is satisfied by ≥1 non-deny match.
    // 3. Otherwise → false.
    const exactMatch = this._matchExact(fullPath);
    const globMatch = this._matchGlob(fullPath);

    // If any deny match exists (exact or glob), short-circuit to false
    if ((exactMatch && this.isDenied(exactMatch.entry)) || (globMatch && this.isDenied(globMatch.entry))) {
      return false;
    }

    // Now check capability satisfaction across the best exact/glob
    const candidates = [exactMatch?.entry, globMatch?.entry].filter(Boolean);

    return capabilities.every((cap) => {
      // Null means “any non-deny capability”
      if (cap === null) return candidates.some((e) => !this.isDenied(e));
      return candidates.some((e) => !this.isDenied(e) && this.hasCapability(e, cap));
    });
  }

  // ===== Public checks =======================================================

  // Entry point used by the has-permissions helper (and tests).
  // If `routeParams` is an array and `requireAll` is true, *all* params must
  // be permitted; otherwise any one is enough.
  // For identity entities/groups pages we require 'list'; other pages only need
  // “not denied”, so we pass `[null]` to mean “any non-deny capability suffices”.
  // Determine if a nav item should be shown.
  // - Resolves one or more API paths from API_PATHS[navItem] (optionally using routeParams).
  // - If routeParams is an array:
  //    • requireAll=true  → every param must be permitted (Array.every)
  //    • requireAll=false → at least one must be permitted (Array.some)
  // - Capability policy:
  //    • 'entities' / 'groups' require ["list"]
  //    • everything else uses [null] → “any non-deny capability is enough”
  // - Delegates to hasPermission(path, capabilities) per resolved path.
  hasNavPermission(navItem, routeParams, requireAll) {
    if (routeParams) {
      const params = Array.isArray(routeParams) ? routeParams : [routeParams];
      const evalMethod = !Array.isArray(routeParams) || requireAll ? 'every' : 'some';
      return params[evalMethod]((param) => {
        const capability = param === 'entities' || param === 'groups' ? ['list'] : [null];
        return this.hasPermission(API_PATHS[navItem][param], capability);
      });
    }
    // No params → show if any canonical path for the nav item is permitted.
    return Object.values(API_PATHS[navItem]).some((path) => this.hasPermission(path));
  }

  // Compute route params (models) for the *first* accessible path of a nav item.
  // - Uses hasPermission(...) to pick the first allowed path in API_PATHS[navItem].
  // - For 'policies' and 'tools', the model is the last path segment (e.g., 'acl', 'hash').
  // - Otherwise, returns the pre-mapped route + models from API_PATHS_TO_ROUTE_PARAMS.
  navPathParams(navItem) {
    const path = Object.values(API_PATHS[navItem]).find((path) => this.hasPermission(path));
    if (['policies', 'tools'].includes(navItem)) {
      const last = path.split('/').pop();
      return { models: [last] };
    }
    return API_PATHS_TO_ROUTE_PARAMS[path];
  }

  // Build a fully-qualified API path scoped to the effective namespace.
  // - Prefixes <chroot>/<currentNamespace> when present.
  // - sanitizePath/sanitizeStart prevent accidental double slashes and leading '/' issues.
  // Examples:
  //   ns="team", chroot=null, path="sys/auth"   → "team/sys/auth"
  //   ns="team/child", chroot="admin", path="/sys/auth/" → "admin/team/child/sys/auth/"
  pathNameWithNamespace(pathName) {
    const ns = this.fullCurrentNamespace;
    return ns ? `${sanitizePath(ns)}/${sanitizeStart(pathName)}` : pathName;
  }

  // Core permission check used by both sidebar gating and banner canary probes.
  // - Root short-circuit: real root (isRoot) → true.
  // - Otherwise, compose the fully-qualified path and delegate to _decide(full, capabilities).
  // - 'capabilities' is the caller’s requirement:
  //     • [null]  → “any non-deny capability is enough” (most nav checks)
  //     • ['list'] for identity entities/groups
  //     • or a concrete set like ['read','update'] for specific flows
  hasPermission(pathName, capabilities = [null]) {
    if (this.isRoot) return true;
    const full = this.pathNameWithNamespace(pathName);
    return this._decide(full, capabilities);
  }

  // ===== capability helpers ==================================================

  // If a specific capability is requested, ensure it’s present.
  // For the “any non-deny” case we pass `[null]` and gate on `isDenied(...)` only.
  hasCapability(path, capability) {
    return path.capabilities.includes(capability);
  }

  // A deny anywhere is authoritative.
  isDenied(path) {
    return path.capabilities.includes('deny');
  }
}
