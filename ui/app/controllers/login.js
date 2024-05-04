import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { storageFor } from 'ember-local-storage';
import { tracked } from '@glimmer/tracking';

export default class LoginController extends Controller {
  @tracked email = '';
  @tracked password = '';
  @tracked error = '';
  @storageFor('config') localConfig;
  @service router;
  @service api;

  @action
  async Login() {
    this.error = '';
    let response = await this.api.Login(this.email, this.password);
    if (!response.ok) {
      this.error = data.message;
      return;
    }
    let data = await response.json();
    this.localConfig.set('user', this.email);
    this.localConfig.set('token', data);
    this.router.transitionTo('index');
  }

  @action
  Register() {}
}
