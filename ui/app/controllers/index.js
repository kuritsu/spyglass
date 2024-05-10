import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { storageFor } from 'ember-local-storage';

export default class IndexController extends Controller {
  @tracked error;
  @tracked textFilter;
  @service componentConfig;
  @storageFor('config') localConfig;

  init(params) {
    super.init(params);
    this.componentConfig.subscribe(this.onConfigChange);
    this.textFilter = this.localConfig.get('textFilter');
  }

  @action
  onConfigChange(prop, value) {
    if (prop == 'fetchError') {
      this.error = value;
      return;
    }
    if (prop == 'textFilter') {
      this.textFilter = value;
      return;
    }
  }
}
