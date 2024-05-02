import Service from '@ember/service';

export default class ComponentConfigService extends Service {
    constructor() {
        super();
        this.config = {
            display: 'ID'
        };
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
        this.callBacks.forEach(function(f) { f(prop, value); });
    }
}
