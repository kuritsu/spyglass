import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class TreeViewFilters extends Component {
  @tracked display = 'ID';

  @action
  updateDisplay(displayParam) {
    this.display = displayParam;
    this.args.onDisplayChange(displayParam);
  }
}
