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
    try {
      let response = await this.api.GetTarget(params.id);
      let data = await response.json();
      if (response.status == 403) {
        // Forbidden
        this.api.LogOut();
        return null;
      }
      if (!response.ok) {
        this.componentConfig.update('fetchError', data.message);
        return [];
      }
      return data;
    } catch (err) {
      this.componentConfig.update('fetchError', (err instanceof TypeError) ? "Network error." : err);
      return [];
    }
  }
}
