import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { storageFor } from 'ember-local-storage';

export default class MainMenu extends Component {
  @tracked display = 'ID';
  @tracked reloadTime = 5000;
  @tracked progressTime = 0;
  @tracked showClass = '';
  @tracked textFilter = '';
  @tracked timeToRefresh = 0;
  @tracked user = '';
  @service componentConfig;
  @service router;
  @service api;
  @storageFor('config') localConfig;

  @action
  init() {
    this.display = this.localConfig.get('display');
    this.reloadTime = this.localConfig.get('reloadTime');
    this.textFilter = this.localConfig.get('textFilter');
    this.user = this.localConfig.get('user');
    this.componentConfig.update('display', this.display);
    this.componentConfig.update('reloadTime', this.reloadTime);
    this.componentConfig.update('textFilter', this.textFilter);
    if (this.reloadTime > 0) {
      setTimeout(this.makeProgress, 1000);
    }
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
    this.timeToRefresh = 0;
    if (this.reloadTime > 0) {
      setTimeout(this.makeProgress, 1000);
    }
    this.localConfig.set('reloadTime', this.reloadTime);
  }

  @action
  makeProgress() {
    this.timeToRefresh += 1000;
    if (this.timeToRefresh == this.reloadTime) {
      this.router.refresh(this.router.currentRouteName);
      this.timeToRefresh = 0;
    }
    setTimeout(this.makeProgress, 1000);
  }

  @action
  filterByText(event) {
    this.updateTextFilter(event.target.value);
  }

  @action
  updateTextFilter(text) {
    this.textFilter = text;
    this.componentConfig.update('textFilter', text);
    this.localConfig.set('textFilter', this.textFilter);
  }

  @action
  LogOut() {
    this.api.LogOut();
  }

  @action
  Settings() {
    this.router.transitionTo('settings');
  }
}
