import Route from '@ember/routing/route';

export default class IndexRoute extends Route {
  model() {
    return [
      {
        id: 'weekdays',
        description: 'Week days',
        status: 30,
        url: 'https://google.com',
        children: [
          { id: 'sunday', description: 'Sunday', status: 0 },
          { id: 'monday', description: 'Monday', status: 30 },
          { id: 'tuesday', description: 'Tuesday', status: 100 },
        ],
      },
    ];
  }
}
