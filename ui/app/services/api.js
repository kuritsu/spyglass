import Service from '@ember/service';
import { service } from '@ember/service';
import { storageFor } from 'ember-local-storage';

export default class ApiService extends Service {
  @storageFor('config') localConfig;
  @service router;

  async Login(email, password) {
    let reqHeaders = new Headers();
    reqHeaders.set('Content-Type', 'application/json');
    let response = await fetch('http://localhost:8010/login', {
      method: 'POST',
      body: JSON.stringify({
        email: email,
        password: password,
      }),
      headers: reqHeaders,
    });
    return response;
  }

  LogOut() {
    this.localConfig.set('user', '');
    this.localConfig.set('token', '');
    this.router.transitionTo('login');
  }

  getToken() {
    let email = this.localConfig.get('user');
    let token = this.localConfig.get('token');
    return `${email}:${token}`;
  }

  createHeaders() {
    let reqHeaders = new Headers();
    reqHeaders.set('Content-Type', 'application/json');
    reqHeaders.set('Authorization', this.getToken());
    return reqHeaders;
  }

  async ListTargets() {
    let response = await fetch('http://localhost:8010/targets', {
      method: 'GET',
      headers: this.createHeaders(),
    });
    return response;
  }

  async GetTarget(id) {
    let response = await fetch(
      `http://localhost:8010/target?id=${id}&includeChildren=true`,
      {
        method: 'GET',
        headers: this.createHeaders(),
      },
    );
    return response;
  }
}
