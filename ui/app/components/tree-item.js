import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class TreeItem extends Component {
  @tracked display = 'ID';
  @service componentConfig;
  @tracked show = '';

  @action
  init() {
    this.componentConfig.subscribe(this.onPropChange);
  }

  @action
  onPropChange(prop, value) {
    if (prop == 'display') {
      this.display = value;
    }
    if (prop == 'textFilter') {
      this.show =
        JSON.stringify(this.args.target)
          .toLowerCase()
          .indexOf(value.toLowerCase()) != -1
          ? ''
          : 'd-none';
    }
  }

  get Style() {
    if (this.args.target.children) {
      for (let i = 0; i < this.args.target.children.Length; i++) {
        if (
          this.args.target.children[i].critical &&
          this.args.target.children[i].status != 100
        ) {
          return 'text-bg-danger';
        }
      }
    }
    if (this.args.target.status == 0) {
      return 'text-bg-danger';
    } else if (this.args.target.status == 100) {
      return 'text-bg-success';
    }
    return 'text-bg-warning';
  }

  get Value() {
    switch (this.display) {
      case 'ID':
        let result = this.args.target.id;
        if (this.args.parent) {
          result = result.substring(this.args.parent.length + 1);
        }
        return result;
      case 'Status':
        return this.args.target.status;
      default:
        return ' ';
    }
  }
}
