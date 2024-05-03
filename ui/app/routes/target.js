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
    return data;
  }
}
