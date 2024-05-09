import Component from '@glimmer/component';

export default class FilteredTreeItems extends Component {
    get results() {
        let { items, query } = this.args;
        if (query) {
            items = items.filter((e) =>
                JSON.stringify(e)
                  .toLowerCase()
                  .indexOf(query.toLowerCase()) != -1)
        }
        return items;
    }
}