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
    this.onPropChange('display', this.componentConfig.get('display'));
    this.onPropChange('textFilter', this.componentConfig.get('textFilter'));
  }

  @action
  onPropChange(prop, value) {
    if (prop == 'display') {
      this.display = value;
    }
  }

  get Style() {
    let append = this.display == 'Status' ? 'statusItem' : '';
    if (this.args.target.children) {
      for (let i = 0; i < this.args.target.children.Length; i++) {
        if (
          this.args.target.children[i].critical &&
          this.args.target.children[i].status != 100
        ) {
          return `treeViewRed ${append}`;
        }
      }
    }
    if (this.args.target.status == 0) {
      return `treeViewRed ${append}`;
    } else if (this.args.target.status == 100) {
      return `treeViewGreen ${append}`;
    }
    return `treeViewYellow ${append}`;
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
        return '.';
    }
  }
}
