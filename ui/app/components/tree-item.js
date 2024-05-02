import Component from '@glimmer/component';

export default class TreeItem extends Component {
  get Style() {
    if (this.args.target.children) {
      for (let i = 0; i < this.args.target.children.Length; i++) {
        if (
          this.args.target.children[i].critical &&
          this.args.target.children[i].status != 100
        ) {
          return 'treeViewRed';
        }
      }
    }
    if (this.args.target.status == 0) {
      return 'treeViewRed';
    } else if (this.args.target.status == 100) {
      return 'treeViewGreen';
    }
    return 'treeViewYellow';
  }

  get Value() {
    switch (this.args.display) {
      case 'ID':
        return this.args.target.id;
      case 'Status':
        return this.args.target.status;
      default:
        return ' ';
    }
  }
}
