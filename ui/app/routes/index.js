import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class IndexRoute extends Route {
  @service componentConfig;
  @service api;

  async model() {
    let response = await this.api.ListTargets();
    let data = await response.json();
    if (response.status == 403) { // Forbidden
      this.api.LogOut();
      return null;
    }
    if (!response.ok) {
      this.componentConfig.update('modelError', data.message);
      return null;
    }
    return data;
  }
}
