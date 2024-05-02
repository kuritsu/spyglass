import Route from '@ember/routing/route';

export default class TargetRoute extends Route {
  queryParams = {
    id: {
      refreshModel: true,
    },
  };

  async model(params) {
    let response = await fetch(
      `http://localhost:8010/target?id=${params['id']}&includeChildren=true`,
    );
    let data = await response.json();
    return data;
  }
}
