import { module, test } from 'qunit';
import pathToTree from 'vault/lib/path-to-tree';

module('Unit | Lib | path to tree', function () {
  const tests = [
    [
      'basic',
      ['one', 'one/two', 'one/two/three/four/five'],
      {
        one: {
          two: {
            three: {
              four: {
                five: null,
              },
            },
          },
        },
      },
    ],
    [
      'multiple leaves on a level',
      ['one', 'two', 'three/four/five', 'three/four/six', 'three/four/six/one'],
      {
        one: null,
        three: {
          four: {
            five: null,
            six: {
              one: null,
            },
          },
        },
        two: null,
      },
    ],
    [
      'leaves with shared prefix',
      ['ns1', 'ns1a', 'ns1a/ns2/ns3', 'ns1a/ns2a/ns3'],
      {
        ns1: null,
        ns1a: {
          ns2: {
            ns3: null,
          },
          ns2a: {
            ns3: null,
          },
        },
      },
    ],
    [
      'leaves with nested number and shared prefix',
      ['ns1', 'ns1a', 'ns1a/99999/five9s', 'ns1a/999/ns3', 'ns1a/9999/ns3'],
      {
        ns1: null,
        ns1a: {
          999: {
            ns3: null,
          },
          9999: {
            ns3: null,
          },
          99999: {
            five9s: null,
          },
        },
      },
    ],
    [
      'sorting lexicographically',
      [
        '99',
        'bat',
        'bat/bird',
        'animal/flying/birds',
        'animal/walking/dogs',
        'animal/walking/cats',
        '1/thing',
      ],
      {
        1: {
          thing: null,
        },
        99: null,
        animal: {
          flying: {
            birds: null,
          },
          walking: {
            cats: null,
            dogs: null,
          },
        },
        bat: {
          bird: null,
        },
      },
    ],
  ];

  tests.forEach(function ([name, input, expected]) {
    test(`pathToTree: ${name}`, function (assert) {
      const output = pathToTree(input);
      assert.deepEqual(output, expected, 'has expected data');
    });
  });
});
