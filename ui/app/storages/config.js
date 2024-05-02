import StorageObject from 'ember-local-storage/local/object';

const Storage = StorageObject.extend();

Storage.reopenClass({
  initialState() {
    return {
        reloadTime: 5000
    };
  }
});

export default Storage;
