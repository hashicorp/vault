import { module, test } from 'qunit';
import pathToTree from 'vault/lib/path-to-tree';

module('Unit | Lib | path to tree', function() {
  let tests = [
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
  ];

  tests.forEach(function([name, input, expected]) {
    test(`pathToTree: ${name}`, function(assert) {
      let output = pathToTree(input);
      assert.deepEqual(output, expected, 'has expected data');
    });
  });
});
