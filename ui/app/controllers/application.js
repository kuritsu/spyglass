import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { storageFor } from 'ember-local-storage';

export default class ApplicationController extends Controller {
  @storageFor('config') localConfig;
  @service router;

  @action
  didInsert() {
    let user = this.localConfig.get('user');
    if (user == '' && this.router.currentRouteName != 'login') {
      this.router.transitionTo('login');
    }
  }

  get Show() {
    return this.router.currentRouteName != 'login';
  }
}
