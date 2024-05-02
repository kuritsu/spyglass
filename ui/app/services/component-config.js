import Service from '@ember/service';

export default class ComponentConfigService extends Service {
  constructor() {
    super();
    this.config = {};
    this.callBacks = [];
  }

  get(prop) {
    return this.config[prop];
  }

  update(prop, value) {
    this.config[prop] = value;
    this.notify(prop, value);
  }

  subscribe(callback) {
    this.callBacks.push(callback);
  }

  notify(prop, value) {
    this.callBacks.forEach(function (f) {
      f(prop, value);
    });
  }

  refreshValues() {
    console.log(this.config);
    keys = Object.keys(this.config);
    values = this.config;
    keys.forEach(function (k) {
      this.callBacks.forEach(function (f) {
        f(k, values[k]);
      });
    });
  }
}
