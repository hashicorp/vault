import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | linkable-item', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders anything passed in', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`<LinkableItem />`);
    assert.equal(this.element.textContent.trim(), '');

    // Template block usage:
    await render(hbs`
      <LinkableItem as |Li|>
        <Li.content>
          stuff here
        </Li.content>
        <Li.Menu>
          menu
        </Li.Menu>
      </LinkableItem>
    `);

    assert.dom('[data-test-linkable-item-content]').hasText('stuff here');
    assert.dom('[data-test-linkable-item-menu]').hasText('menu');
  });

  test('it is not wrapped in a linked block if no link is passed', async function(assert) {
    await render(hbs`
      <LinkableItem as |Li|>
        <Li.content>
          stuff here
        </Li.content>
      </LinkableItem>
    `);

    assert.dom('.list-item-row').exists('List item row exists');
    assert.dom('.list-item-row.linked-block').doesNotExist('Does not render linked block');
  });
  test('it is wrapped in a linked block if a link is passed', async function(assert) {
    await render(hbs`
      <LinkableItem @link={{hash route="vault" model="modelId"}} as |Li|>
        <Li.content
          @title={{title}}
          @link={{hash route="vault" model="modelId"}}
        >
          stuff here
        </Li.content>
      </LinkableItem>
    `);

    assert.dom('.list-item-row.linked-block').exists('Renders linked block');
  });

  test('it renders standard attributes on content', async function(assert) {
    this.set('title', 'Hello');
    this.set('accessor', 'my accessor');
    this.set('description', 'my description');
    this.set('glyph', 'key');
    this.set('glyphText', 'Here is some extra info');

    // Template block usage:
    await render(hbs`
      <LinkableItem data-test-example as |Li|>
        <Li.content
          @accessor="my accessor"
          @description="Description goes here
          @glyphText={{backend.engineType}}
          @glyph={{or (if (eq backend.engineType "kmip") "secrets" backend.engineType) "secrets"}}
          @title={{title}} 
        >
          stuff here
        </Li.content>
      </LinkableItem>
    `);

    assert.dom('.list-item-row').exists('List item row exists');
    assert.dom('.list-item-row.linked-block').doesNotExist('Does not render linked block');
  });

  // Optional case
  // test('it shows title as glyph tooltip if no glyphtext provided', async function(assert) {
  // });
});
