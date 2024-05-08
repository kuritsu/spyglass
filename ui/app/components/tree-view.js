import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class TreeView extends Component {
    @tracked modalOpen = false;
    @tracked childrenVisibleCount = 0;
    @service componentConfig;
    @tracked filteredChildren;

    constructor() {
        super(...arguments);
        this.filteredChildren = this.args.target.children;
        this.childrenVisibleCount = this.filteredChildren ? this.filteredChildren.length : 0;
        this.componentConfig.subscribe(this.onPropChange);
    }

    @action
    onPropChange(prop, value) {
        if (prop != 'textFilter' || !this.args.target.children)
            return;
        let selected = [];
        this.args.target.children.forEach(e => {
            if (JSON.stringify(e).toLowerCase().indexOf(value.toLowerCase()) != -1)
                selected.push(e);
        });
        this.filteredChildren = selected;
        this.childrenVisibleCount = this.filteredChildren.length;
    }

    @action
    ShowModal() {
        this.modalOpen = true;
    }

    @action
    getChildId(id) {
        let lastSlash = id.lastIndexOf('/');
        return lastSlash > -1 ? id.substring(lastSlash + 1) : id;
    }
}