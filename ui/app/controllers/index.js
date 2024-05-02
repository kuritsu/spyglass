import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class IndexController extends Controller {
  @tracked display = 'ID';

  @action
  onDisplayChange(display) {
    this.display = display;
  }
}
