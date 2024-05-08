import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class TreeView extends Component {
    @tracked
    modalOpen = false

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