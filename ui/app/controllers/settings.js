import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { storageFor } from 'ember-local-storage';

export default class ChangePasswordController extends Controller {
  @tracked oldPassword;
  @tracked newPassword;
  @tracked confirmPassword;
  @tracked error;
  @service api;
  @service router;
  @storageFor('config') localConfig;

  @action
  async Update(event) {
    this.error = '';
    if (!document.forms['passwordChange'].reportValidity()) return;
    event.preventDefault();
    if (this.newPassword != this.confirmPassword) {
      this.error = "Passwords don't match.";
      return;
    }
    event.target.disabled = true;
    try {
      let response = await this.api.UpdateUser(
        '',
        this.oldPassword,
        this.newPassword,
      );
      let data = await response.json();
      if (!response.ok) {
        this.error = data.message;
      } else {
        this.router.transitionTo('index');
        return;
      }
    } catch (error) {
      console.log(error);
      this.error = 'Network error';
    }
    event.target.disabled = false;
  }
}
