import StorageObject from 'ember-local-storage/local/object';

const Storage = StorageObject.extend();

Storage.reopenClass({
  initialState() {
    return {
        display: 'ID',
        reloadTime: 5000,
        textFilter: ''
    };
  }
});

export default Storage;
