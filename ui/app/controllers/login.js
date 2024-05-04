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
  async Login(event) {
    event.preventDefault();
    event.target.disabled = true;
    this.error = '';
    try {
      let response = await this.api.Login(this.email, this.password);
      if (!response.ok) {
        this.error = data.message;
        event.target.disabled = false;
        return;
      }
      let data = await response.json();
      this.localConfig.set('user', this.email);
      this.localConfig.set('token', data);
      this.router.transitionTo('index');
    } catch (error) {
      this.error = 'Network error';
      event.target.disabled = false;
    }
  }

  @action
  Register() {}
}
