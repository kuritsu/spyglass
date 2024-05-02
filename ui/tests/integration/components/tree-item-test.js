import { module, test } from 'qunit';
import { setupRenderingTest } from 'ui/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | tree-item', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`<TreeItem />`);

    assert.dom().hasText('');

    // Template block usage:
    await render(hbs`
      <TreeItem>
        template block text
      </TreeItem>
    `);

    assert.dom().hasText('template block text');
  });
});
