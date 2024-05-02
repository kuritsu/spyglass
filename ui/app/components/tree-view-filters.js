import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';

export default class TreeViewFilters extends Component {
  @tracked display = 'ID';
  @tracked reloadTime = 5000;
  @tracked showClass = "";
  @service componentConfig;

  @action
  updateDisplay(displayParam) {
    this.display = displayParam;
    this.componentConfig.update('display', displayParam);
  }

  get ReloadTime() {
    if (this.reloadTime === 0) {
        return `Never`;
    }
    if (this.reloadTime < 60000) { // 1 minute 
        return `Every ${this.reloadTime / 1000}s`;
    }
    return `Every ${this.reloadTime / 60000}m`;
  }

  @action
  updateReload(reloadParam) {
    this.reloadTime = reloadParam * 1000;
  }

  @action
  hide() {
    this.showClass = "d-none";
  }
}
