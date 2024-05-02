import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { storageFor } from 'ember-local-storage';

export default class TreeViewFilters extends Component {
  @tracked display = 'ID';
  @tracked reloadTime = 5000;
  @tracked showClass = '';
  @tracked textFilter = '';
  @service componentConfig;
  @storageFor('config') localConfig;

  @action
  init() {
    this.updateReload(this.localConfig.get('reloadTime') / 1000);
  }

  @action
  updateDisplay(displayParam) {
    this.display = displayParam;
    this.componentConfig.update('display', displayParam);
  }

  get ReloadTime() {
    if (this.reloadTime === 0) {
      return `Never`;
    }
    if (this.reloadTime < 60000) {
      // 1 minute
      return `Every ${this.reloadTime / 1000}s`;
    }
    return `Every ${this.reloadTime / 60000}m`;
  }

  @action
  updateReload(reloadParam) {
    this.reloadTime = reloadParam * 1000;
    if (this.reloadTime > 0) {
      setTimeout(this.refresh, this.reloadTime);
    }
    this.localConfig.set('reloadTime', this.reloadTime);
  }

  @action
  refresh() {
    location.reload();
    if (this.reloadTime > 0) {
      setTimeout(this.refresh, this.reloadTime);
    }
  }

  @action
  hide() {
    this.showClass = 'd-none';
  }

  @action
  filterByText(event) {
    this.textFilter = event.target.value;
    this.componentConfig.update('textFilter', event.target.value);
  }
}
