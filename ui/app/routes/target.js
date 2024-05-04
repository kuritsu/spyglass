import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class TargetRoute extends Route {
  @service api;

  queryParams = {
    id: {
      refreshModel: true,
    },
  };

  async model(params) {
    let response = await this.api.GetTarget(params.id);
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
