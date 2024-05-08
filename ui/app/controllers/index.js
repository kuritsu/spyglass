import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class IndexController extends Controller {
  @tracked error;
  @tracked filteredTargets = null;
  @service componentConfig;

  init(params) {
    super.init(params);
    this.componentConfig.subscribe(this.onConfigChange);
  }

  @action
  onConfigChange(prop, value) {
    if (prop == 'fetchError') {
      this.error = value;
      return;
    }
    if (prop == 'textFilter') {
        if (!this.model)
          return;
        let selected = [];
        this.model.forEach(e => {
            if (JSON.stringify(e).toLowerCase().indexOf(value.toLowerCase()) != -1)
                selected.push(e);
        });
        this.filteredTargets = selected;
    }
  }
}
