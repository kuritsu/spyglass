import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class TargetRoute extends Route {
  @service api;
  @service componentConfig;

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
        return null;
      }
      //TODO: Retrieve parent info from api
      let sepIndex = params.id.lastIndexOf('/');
      if (sepIndex > -1) {
        data.parent = {
          id: params.id.substring(0, sepIndex),
          status: 0,
        };
      }
      return data;
    } catch (err) {
      this.componentConfig.update(
        'fetchError',
        err instanceof TypeError ? 'Network error.' : err,
      );
      return null;
    }
  }
}
