import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class IndexRoute extends Route {
  @service api;

  async model() {
    let response = await this.api.ListTargets();
    let data = await response.json();
    return data;
  }
}
