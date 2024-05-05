import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class IndexController extends Controller {
  @tracked error;
  @service componentConfig;

  init(params) {
    super.init(params);
    this.componentConfig.subscribe(this.onConfigChange);
  }

  @action
  onConfigChange(prop, value) {
    if (prop != 'fetchError') {
      return;
    }
    this.error = value;
  }
}
