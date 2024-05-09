import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { storageFor } from 'ember-local-storage';

export default class TreeView extends Component {
  @tracked modalOpen = false;
  @tracked childrenVisibleCount = 0;
  @service componentConfig;
  @tracked filteredChildren;
  @tracked textFilter = '';
  @storageFor('config') localConfig;

  constructor(...args) {
    super(...args);
    this.componentConfig.subscribe(this.onPropChange);
    this.textFilter = this.localConfig.get('textFilter');
  }

  @action
  onPropChange(prop, value) {
    if (prop != 'textFilter') return;
    this.textFilter = value;
  }

  @action
  ShowModal() {
    this.modalOpen = true;
  }

  @action
  getChildId(id) {
    let lastSlash = id.lastIndexOf('/');
    return lastSlash > -1 ? id.substring(lastSlash + 1) : id;
  }
}
