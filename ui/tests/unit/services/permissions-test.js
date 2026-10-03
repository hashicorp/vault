/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import Service from '@ember/service';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { PERMISSIONS_BANNER_STATES, RESULTANT_ACL_PATH } from 'vault/services/permissions';

const PERMISSIONS_RESPONSE = {
  data: {
    exact_paths: {
      foo: { capabilities: ['read'] },
      'bar/bee': { capabilities: ['create', 'list'] },
      boo: { capabilities: ['deny'] },
    },
    glob_paths: {
      'baz/biz': { capabilities: ['read'] },
      'ends/in/slash/': { capabilities: ['list'] },
    },
  },
};

// Small helper to DRY namespace registration across modules
function registerNs(owner, path) {
  const Ns = Service.extend({ path });
  owner.register('service:namespace', Ns);
}

// Note: Policy matching follows Vault’s priority rules
// (see: https://developer.hashicorp.com/vault/docs/concepts/policies#priority-matching).
// In short:
// - '+' = exactly one segment, '*' = wildcard, '/' = include base, '/*' = children-only.
// - Most-specific (longest) match wins; deny beats allow.
// - Exact and glob matches both apply; any deny blocks access.
// - Special: '' key means allow-all (`path "*" { ... }`), but can still be overridden by a more specific deny.
// - Root tokens short-circuit to allow; banners use canary-path heuristics.

module('Unit | Service | permissions', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.get(`/${RESULTANT_ACL_PATH}`, () => PERMISSIONS_RESPONSE);
    this.service = this.owner.lookup('service:permissions');
  });

  // ───────────────────────────────────────────────────────────────────────────
  // Basics
  // ───────────────────────────────────────────────────────────────────────────
  test('sets paths properly', async function (assert) {
    await this.service.getPaths.perform();
    assert.deepEqual(this.service.exactPaths, PERMISSIONS_RESPONSE.data.exact_paths);
    assert.deepEqual(this.service.globPaths, PERMISSIONS_RESPONSE.data.glob_paths);
  });

  test('sets the root token', function (assert) {
    this.service.setPaths({ data: { root: true } });
    assert.true(this.service.isRoot, 'isRoot is true with root token');
  });

  test('defaults to show all items when policy cannot be found', async function (assert) {
    this.server.get(`${RESULTANT_ACL_PATH}`, () => overrideResponse(403));
    await this.service.getPaths.perform();
    assert.true(this.service.hasFallbackAccess, 'hasFallbackAccess is true when policy cannot be found');
  });

  // ───────────────────────────────────────────────────────────────────────────
  // navPathParams
  // ───────────────────────────────────────────────────────────────────────────
  test('returns the first allowed nav route for policies', function (assert) {
    const policyPaths = {
      'sys/policies/acl': { capabilities: ['deny'] },
      'sys/policies/rgp': { capabilities: ['read'] },
    };
    this.service.setProperties({ exactPaths: policyPaths });
    assert.strictEqual(
      this.service.navPathParams('policies').models[0],
      'rgp',
      'first allowed route is returned'
    );
  });

  test('returns the first allowed nav route for access', function (assert) {
    const accessPaths = {
      'sys/auth': { capabilities: ['deny'] },
      'identity/entity/id': { capabilities: ['read'] },
    };
    const expected = { route: 'vault.cluster.access.identity', models: ['entities'] };
    this.service.setProperties({ exactPaths: accessPaths });
    assert.deepEqual(this.service.navPathParams('access'), expected);
  });

  test('navPathParams uses canonical first route when a glob allows all', function (assert) {
    this.service.setProperties({ globPaths: { 'sys/policies/*': { capabilities: ['read'] } } });
    const result = this.service.navPathParams('policies');
    assert.deepEqual(
      result,
      { models: ['acl'] },
      'picks first in canonical order (acl → rgp → egp) when all allowed'
    );
  });

  test('most-specific exact key wins (length-descending)', function (assert) {
    this.service.setProperties({
      exactPaths: {
        'sys/auth': { capabilities: ['deny'] },
        'sys/auth/methods': { capabilities: ['read'] },
      },
    });
    // Querying base "sys/auth" should still be allowed if a deeper exact key is intended to satisfy children,
    // BUT our matcher selects the most specific only when the base is equal or a parent. For "sys/auth/methods"
    // we expect the deeper allow to win; for the base itself we respect the base entry.
    assert.true(this.service.hasPermission('sys/auth/methods'), 'deeper exact allow applies');
    assert.false(this.service.hasPermission('sys/auth'), 'base respects base deny');
  });

  module('hasPermission', function () {
    test('returns true if a policy includes access to an exact path', function (assert) {
      this.service.setProperties({ exactPaths: PERMISSIONS_RESPONSE.data.exact_paths });
      assert.true(this.service.hasPermission('foo'), 'policy includes access to foo exact path');
    });

    test("returns true if a path's base is included in the policy exact paths", function (assert) {
      this.service.setProperties({ exactPaths: PERMISSIONS_RESPONSE.data.exact_paths });
      assert.true(this.service.hasPermission('bar'), 'base "bar" satisfied by "bar/bee" key');
    });

    // An empty string key here represents a policy of `path "*" { ... }`
    // (wildcard for all paths). While it can be argued this doesn't guarantee meaningful access to canary endpoints,
    // we are interpreting it to mean "allow all."
    // This is consistent with the behavior of the root token, which also suppresses the banner.
    test('empty-root glob key ("") (e.g. sys "*") interrupted as allow all', function (assert) {
      this.service.setProperties({ globPaths: { '': { capabilities: ['read'] } } });
      assert.true(this.service.hasPermission('sys/auth'), 'empty-root allows all');
    });

    test('most-specific deny overrides empty-root allow', function (assert) {
      // "" = allow-all, but a longer, more specific deny should win
      this.service.setProperties({
        globPaths: {
          '': { capabilities: ['read'] },
          'sys/auth/*': { capabilities: ['deny'] },
        },
      });
      assert.false(
        this.service.hasPermission('sys/auth/methods'),
        'specific deny on sys/auth/* wins over allow-all'
      );
      assert.true(
        this.service.hasPermission('sys/policies/acl'),
        'non-matching path remains allowed via empty-root allow'
      );
    });

    test('returns false if the matched exact path includes deny capability', function (assert) {
      this.service.setProperties({ exactPaths: { boo: { capabilities: ['deny'] } } });
      assert.false(this.service.hasPermission('boo'), 'deny capability blocks access');
    });

    test('matches whether or not path ends with slash when glob ends with slash', function (assert) {
      this.service.setProperties({ globPaths: { 'ends/in/slash/': { capabilities: ['list'] } } });
      assert.true(this.service.hasPermission('ends/in/slash'), 'matches without slash');
      assert.true(this.service.hasPermission('ends/in/slash/'), 'matches with slash');
    });

    test('returns false if a policy does not include access to a path', function (assert) {
      assert.false(this.service.hasPermission('danger'));
    });

    test('returns true with the root token', function (assert) {
      this.service.setProperties({ isRoot: true });
      assert.true(this.service.hasPermission('hi'));
    });

    test('returns true if policy has all requested capabilities on a path', function (assert) {
      this.service.setProperties({
        exactPaths: PERMISSIONS_RESPONSE.data.exact_paths,
        globPaths: PERMISSIONS_RESPONSE.data.glob_paths,
      });
      assert.true(this.service.hasPermission('bar/bee', ['create', 'list']));
      assert.true(this.service.hasPermission('baz/biz', ['read']));
    });

    test('returns false if policy lacks any requested capability on a path', function (assert) {
      this.service.setProperties({
        exactPaths: PERMISSIONS_RESPONSE.data.exact_paths,
        globPaths: PERMISSIONS_RESPONSE.data.glob_paths,
      });
      assert.false(this.service.hasPermission('bar/bee', ['create', 'delete']), 'delete missing');
      assert.false(this.service.hasPermission('foo', ['create']), 'create missing');
    });

    test('returns false when an exact path matches but the requested capability is missing', function (assert) {
      this.service.setProperties({ exactPaths: { 'bar/bee': { capabilities: ['list'] } } });
      assert.false(this.service.hasPermission('bar/bee', ['update']), 'update not granted');
      assert.false(
        this.service.hasPermission('bar/bee', ['list', 'create']),
        'require-all semantics: create not granted'
      );
    });
  });

  module('hasNavPermission', function () {
    test('returns true if a policy includes the required capabilities for at least one path', function (assert) {
      const accessPaths = {
        'sys/auth': { capabilities: ['deny'] },
        'identity/group/id': { capabilities: ['list', 'read'] },
      };
      this.service.setProperties({ exactPaths: accessPaths });
      assert.true(this.service.hasNavPermission('access', 'groups'));
    });

    test('returns false if a policy lacks the required capabilities for the path', function (assert) {
      const accessPaths = {
        'sys/auth': { capabilities: ['deny'] },
        'identity/group/id': { capabilities: ['read'] },
      };
      this.service.setProperties({ exactPaths: accessPaths });
      assert.false(this.service.hasNavPermission('access', 'groups'));
    });

    test('handles routeParams as array with requireAll semantics', function (assert) {
      const getPaths = (override) => ({
        'sys/auth': { capabilities: [override || 'read'] },
        'identity/mfa/method': { capabilities: [override || 'read'] },
        'identity/oidc/client': { capabilities: [override || 'deny'] },
      });

      this.service.setProperties({ exactPaths: getPaths() });
      assert.true(
        this.service.hasNavPermission('access', ['methods', 'mfa', 'oidc']),
        'true when any route is permitted'
      );
      assert.false(
        this.service.hasNavPermission('access', ['methods', 'mfa', 'oidc'], true),
        'false when any route is not permitted and requireAll is passed'
      );

      this.service.setProperties({ exactPaths: getPaths('read') });
      assert.true(
        this.service.hasNavPermission('access', ['methods', 'mfa', 'oidc'], true),
        'true when all routes are permitted and requireAll is passed'
      );

      this.service.setProperties({ exactPaths: getPaths('deny') });
      assert.false(
        this.service.hasNavPermission('access', ['methods', 'mfa', 'oidc']),
        'false when no routes are permitted'
      );
      assert.false(
        this.service.hasNavPermission('access', ['methods', 'mfa', 'oidc'], true),
        'false when no routes are permitted and requireAll is passed'
      );
    });
  });

  module('pathNameWithNamespace', function () {
    test('appends the namespace to the path if there is one', function (assert) {
      registerNs(this.owner, 'marketing');
      assert.strictEqual(this.service.pathNameWithNamespace('sys/auth'), 'marketing/sys/auth');
    });

    test('appends the chroot and namespace when both present', function (assert) {
      registerNs(this.owner, 'marketing');
      this.service.setProperties({ chrootNamespace: 'admin/' });
      assert.strictEqual(this.service.pathNameWithNamespace('sys/auth'), 'admin/marketing/sys/auth');
    });

    test('appends the chroot when no namespace', function (assert) {
      this.service.setProperties({ chrootNamespace: 'admin' });
      assert.strictEqual(this.service.pathNameWithNamespace('sys/auth'), 'admin/sys/auth');
    });

    test('handles superfluous slashes', function (assert) {
      registerNs(this.owner, '/marketing');
      this.service.setProperties({ chrootNamespace: '/admin/' });
      assert.strictEqual(this.service.pathNameWithNamespace('/sys/auth'), 'admin/marketing/sys/auth');
      assert.strictEqual(
        this.service.pathNameWithNamespace('/sys/policies/'),
        'admin/marketing/sys/policies/'
      );
    });
  });

  // ───────────────────────────────────────────────────────────────────────────
  // Glob matching semantics (+ and *)
  // ───────────────────────────────────────────────────────────────────────────
  module('glob matching semantics', function () {
    test('+ matches exactly one segment; trailing /* means any child (≥1)', function (assert) {
      this.service.setProperties({ globPaths: { 'secret/+/*': { capabilities: ['read'] } } });

      assert.true(
        this.service.hasPermission('secret/team-a/foo'),
        'one segment after secret, then a child → match'
      );
      assert.true(this.service.hasPermission('secret/team-a/foo/bar'), 'multiple children via /* → match');

      // children-only: base MUST have at least one child
      assert.false(
        this.service.hasPermission('secret/foo'),
        'zero segments after + (needs child due to /*) → no match'
      );
      assert.false(this.service.hasPermission('secret/team-a'), 'base-only (no child) with /* → no match');

      // IMPORTANT CHANGE: with + matching "team-a", the /* tail matches any depth of children
      assert.true(
        this.service.hasPermission('secret/team-a/eng/foo/bar/baz'),
        '/* tail matches one-or-more segments; deep children → match'
      );
    });

    test('* greedy tail works', function (assert) {
      this.service.set('globPaths', { 'sys/*': { capabilities: ['list'] } });
      assert.true(this.service.hasPermission('sys/policies/acl'));
      assert.true(this.service.hasPermission('sys/auth/methods'));
    });

    test('include-base semantics for keys ending with "/" (base OR children)', function (assert) {
      this.service.set('globPaths', { 'sys/': { capabilities: ['list'] } });
      assert.true(this.service.hasPermission('sys'), 'base matches');
      assert.true(this.service.hasPermission('sys/'), 'normalized base matches');
      assert.true(this.service.hasPermission('sys/auth'), 'child matches');
    });

    test('most-specific wins across multiple matching globs', function (assert) {
      this.service.set('globPaths', {
        'sys/*': { capabilities: ['deny'] }, // broad deny
        'sys/auth/*': { capabilities: ['read'] }, // specific allow
      });
      assert.true(this.service.hasPermission('sys/auth/methods'), 'specific allow overrides broader deny');
      assert.false(this.service.hasPermission('sys/policies/acl'), 'non-auth still denied by broad rule');
    });

    test('most-specific deny takes precedence when the most specific entry is deny', function (assert) {
      this.service.set('globPaths', {
        'sys/*': { capabilities: ['list'] }, // broad allow
        'sys/auth/*': { capabilities: ['deny'] }, // specific deny
      });
      assert.false(this.service.hasPermission('sys/auth/methods'), 'specific deny blocks');
      assert.true(this.service.hasPermission('sys/policies/acl'), 'non-auth remains allowed by broad rule');
    });

    test('capability filtering applies to glob matches', function (assert) {
      this.service.set('globPaths', { 'sys/tools/*': { capabilities: ['update'] } });
      assert.true(this.service.hasPermission('sys/tools/hash', ['update']));
      assert.false(this.service.hasPermission('sys/tools/hash', ['read']));
    });

    test('deny precedence: exact deny beats glob allow', function (assert) {
      // /ns/sys/auth matches BOTH exact (deny) and glob (allow) → expect false
      this.service.setProperties({
        exactPaths: { 'sys/auth': { capabilities: ['deny'] } },
        globPaths: { 'sys/*': { capabilities: ['read'] } },
      });
      assert.false(this.service.hasPermission('sys/auth'), 'exact deny wins');
    });

    test('deny precedence: glob deny beats exact allow', function (assert) {
      this.service.setProperties({
        exactPaths: { 'sys/auth': { capabilities: ['read'] } },
        globPaths: { 'sys/auth/*': { capabilities: ['deny'] } },
      });
      assert.false(this.service.hasPermission('sys/auth/methods'), 'glob deny wins');
    });
  });

  module('permissionsBanner - fullCurrentNamespace (e.g. ns/<path>)', function () {
    [
      // First set: no chroot or user root
      {
        scenario: 'no chroot or user root → full ns is currentNs',
        chroot: null,
        userRoot: '',
        currentNs: 'foo/bar',
        expectedFullNs: 'foo/bar',
      },
      // Second set: chroot and user root (currentNs excludes chroot)
      {
        scenario: 'chroot + user root → both prefixed',
        chroot: 'foo/',
        userRoot: 'bar',
        currentNs: 'bar/baz',
        expectedFullNs: 'foo/bar/baz',
      },
      // Third set: chroot only
      {
        scenario: 'chroot only → chroot prefixed',
        chroot: 'admin/',
        userRoot: '',
        currentNs: 'child',
        expectedFullNs: 'admin/child',
      },
      // Fourth set: user root only
      {
        scenario: 'user root only → user root prefixed',
        chroot: null,
        userRoot: 'foo',
        currentNs: 'foo/bing',
        expectedFullNs: 'foo/bing',
      },
    ].forEach((testCase) => {
      test(testCase.scenario, async function (assert) {
        const Ns = Service.extend({ userRootNamespace: testCase.userRoot, path: testCase.currentNs });
        this.owner.register('service:namespace', Ns);

        this.service.setPaths({
          data: { glob_paths: {}, exact_paths: {}, chroot_namespace: testCase.chroot },
        });

        const fullNamespace = this.service.fullCurrentNamespace;
        assert.strictEqual(fullNamespace, testCase.expectedFullNs, 'full ns computed correctly');
      });
    });
  });

  module('permissionsBanner — namespace + globPaths coverage (e.g. +/<path>)', function () {
    test('ACL not loaded → no banner (prevents flicker)', function (assert) {
      registerNs(this.owner, 'ns1/child');
      // Fresh service: exact/glob/canViewAll all null
      assert.strictEqual(this.service.permissionsBanner, null, 'null while loading');
    });

    test('ACL load failed → read-failed banner', function (assert) {
      registerNs(this.owner, 'ns1');
      this.service.setProperties({
        _aclLoadFailed: true,
        exactPaths: null,
        globPaths: null,
        canViewAll: false,
      });
      assert.strictEqual(this.service.permissionsBanner, PERMISSIONS_BANNER_STATES.readFailed);
    });

    // If you want root to *always* suppress banners, flip getter order and update this.
    test('root token → no banner regardless of namespace', function (assert) {
      registerNs(this.owner, 'ns1/child');
      this.service.setProperties({
        isRoot: true,
        exactPaths: {},
        globPaths: {},
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, null);
    });

    test('no canaries allowed in ns → no-ns-access banner', function (assert) {
      registerNs(this.owner, 'ns1');
      this.service.setProperties({
        exactPaths: { 'some/noncanary': { capabilities: ['read'] } }, // non-canary
        globPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, PERMISSIONS_BANNER_STATES.noAccess);
    });

    // +/sys/auth (one ns segment) vs current ns depth 1
    test('+/sys/auth (base-only) grants canary in one-segment ns (base only)', function (assert) {
      registerNs(this.owner, 'ns1'); // full: ns1/sys/auth
      this.service.setProperties({
        globPaths: { '+/sys/auth': { capabilities: ['list'] } }, // no trailing slash → base only
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, null, 'base canary suppresses banner');
    });

    test('+/sys/auth/ (include-base) grants canary for base and children in one-segment ns', function (assert) {
      registerNs(this.owner, 'ns1');
      this.service.setProperties({
        globPaths: { '+/sys/auth/': { capabilities: ['list'] } }, // include-base
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, null);
    });

    // +/+/sys/auth for two-segment namespace
    test('+/+/sys/auth (base-only) grants canary in two-segment ns', function (assert) {
      registerNs(this.owner, 'ns1/child'); // full: ns1/child/sys/auth
      this.service.setProperties({
        globPaths: { '+/+/sys/auth': { capabilities: ['read'] } },
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, null);
    });

    test('single + does NOT match two-segment ns → no-ns-access', function (assert) {
      registerNs(this.owner, 'ns1/child');
      this.service.setProperties({
        globPaths: { '+/sys/auth': { capabilities: ['read'] } }, // only one '+'
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, PERMISSIONS_BANNER_STATES.noAccess);
    });

    test('+/+/sys/auth/ include-base still suppresses banner in two-segment ns', function (assert) {
      registerNs(this.owner, 'ns1/child');
      this.service.setProperties({
        globPaths: { '+/+/sys/auth/': { capabilities: ['read'] } },
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, null);
    });

    test('chroot + one-segment ns matches +/+ canary', function (assert) {
      registerNs(this.owner, 'marketing'); // full: admin/marketing/sys/auth
      this.service.setProperties({
        chrootNamespace: 'admin',
        globPaths: { '+/+/sys/auth': { capabilities: ['list'] } },
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, null);
    });

    test('most-specific allow (+/+/sys/auth/) overrides broader deny (+/+/sys/*)', function (assert) {
      registerNs(this.owner, 'ns1/child'); // full path prefix: ns1/child

      this.service.setProperties({
        globPaths: {
          '+/+/sys/*': { capabilities: ['deny'] }, // children-only deny under sys
          '+/+/sys/auth/': { capabilities: ['read'] }, // include-base allow for sys/auth
        },
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });

      assert.strictEqual(
        this.service.permissionsBanner,
        null,
        'specific include-base allow suppresses banner even with broader children-only deny'
      );
    });

    test('deny on the canary path blocks banner suppression when no other canaries are allowed', function (assert) {
      registerNs(this.owner, 'ns1/child');

      this.service.setProperties({
        globPaths: {
          '+/+/sys/auth/': { capabilities: ['deny'] }, // deny the canary
          // no other canary allows present
        },
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });

      assert.strictEqual(
        this.service.permissionsBanner,
        PERMISSIONS_BANNER_STATES.noAccess,
        'without other allowed canaries, deny on sys/auth shows the banner'
      );
    });

    test('empty-root glob "" (aka `path "*"`) suppresses banner in namespaced paths', function (assert) {
      registerNs(this.owner, 'ns1');
      this.service.setProperties({
        globPaths: { '': { capabilities: ['read'] } }, // allow-all
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(
        this.service.permissionsBanner,
        null,
        'allow-all policy suppresses the no-access banner even in a namespace'
      );
    });

    test('canary denied explicitly → no-ns-access', function (assert) {
      registerNs(this.owner, 'ns1/child');
      this.service.setProperties({
        globPaths: { '+/+/sys/auth': { capabilities: ['deny'] } },
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, PERMISSIONS_BANNER_STATES.noAccess);
    });

    test('canary allowed with any non-deny capability suppresses banner', function (assert) {
      registerNs(this.owner, 'ns1/child');
      this.service.setProperties({
        globPaths: { '+/+/sys/auth': { capabilities: ['list'] } }, // any non-deny is enough
        exactPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, null);
    });
  });

  module('permissionsBanner — canary checks (aka “standard mgmt endpoints”)', function () {
    /**
     * Canary endpoints are a curated set of safe, management-style API paths that the UI probes
     * inside the *current* namespace to decide whether to suppress the “no access” banner.
     * Heuristic: any non-deny capability on at least one canary → suppress banner.
     * Examples (not exhaustive): 'sys/auth', 'identity/*', 'sys/leases/lookup', 'sys/policies/*',
     * 'sys/tools/hash', 'sys/replication', 'sys/license', etc.
     */

    test('no banner when any canary path is allowed via exact path (e.g., sys/tools/hash)', function (assert) {
      registerNs(this.owner, '');
      this.service.setProperties({
        exactPaths: { 'sys/tools/hash': { capabilities: ['update'] } }, // canary
        globPaths: {},
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, null, 'canary access suppresses banner');
    });

    test('edge case: namespace-only engine access (+/kv/data/*) still shows banner', function (assert) {
      registerNs(this.owner, 'ns1');
      this.service.setProperties({
        exactPaths: {},
        globPaths: { '+/kv/data/*': { capabilities: ['read'] } }, // non-canary
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(
        this.service.permissionsBanner,
        PERMISSIONS_BANNER_STATES.noAccess,
        'engine-only glob (+/kv/data/*) is not a canary → banner shows'
      );
    });

    test('canary via glob (+/sys/auth/) suppresses banner in a one-segment namespace', function (assert) {
      registerNs(this.owner, 'ns1');
      this.service.setProperties({
        exactPaths: {},
        globPaths: { '+/sys/auth/': { capabilities: ['list'] } }, // include-base glob canary
        canViewAll: false,
        _aclLoadFailed: false,
      });
      assert.strictEqual(this.service.permissionsBanner, null, 'glob canary suppresses banner');
    });
  });

  module('permissionsBanner — resultant-acl direct check', function () {
    test('allow on resultant-acl (ns via + glob) suppresses banner', function (assert) {
      // current namespace depth = 1 → use single '+'
      registerNs(this.owner, 'ns1');
      this.service.setProperties({
        exactPaths: {},
        globPaths: { [`+/${RESULTANT_ACL_PATH}`]: { capabilities: ['read'] } },
        _aclLoadFailed: false,
      });
      assert.strictEqual(
        this.service.permissionsBanner,
        null,
        'read capability to resultant-acl suppresses banner'
      );
    });

    test('deny on resultant-acl (two-segment ns) with no other canaries → no-ns-access', function (assert) {
      registerNs(this.owner, 'ns1/child'); // full path checked: ns1/child/sys/internal/ui/resultant-acl
      this.service.setProperties({
        exactPaths: {}, // no exact allows
        globPaths: {
          [`+/+/${RESULTANT_ACL_PATH}`]: { capabilities: ['deny'] }, // explicit deny
        },
        _aclLoadFailed: false,
      });
      assert.strictEqual(
        this.service.permissionsBanner,
        PERMISSIONS_BANNER_STATES.noAccess,
        'deny to resultant-acl with no other canaries shows no-access banner'
      );
    });
  });
  // ───────────────────────────────────────────────────────────────────────────
  // precedence test (document current getter order)
  // ───────────────────────────────────────────────────────────────────────────
  test('root overrides read-failed: no banner for canViewAll', function (assert) {
    // namespace doesn’t matter; root short-circuits
    const Ns = Service.extend({ path: '' });
    this.owner.register('service:namespace', Ns);

    this.service.setProperties({
      _aclLoadFailed: true, // simulate fetch failure
      isRoot: true,
      canViewAll: true,
      exactPaths: null,
      globPaths: null,
    });

    assert.strictEqual(this.service.permissionsBanner, null, 'root sees no banner even if ACL fetch failed');
  });
});
