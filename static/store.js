import Reflux from "reflux";
// https://github.com/visionmedia/superagent/wiki/Superagent-for-Webpack
var request = require('superagent');

import Actions from "./actions";

var Store = Reflux.createStore({
    listenables: [Actions],

    init: function() {
        this.listenTo(Actions.load, this.fetchData);
    },

    getInitialState() {
        this.list = [];
        return this.list;
    },

    fetchData: function() {
        request
        .get('http://localhost:8000/torrents.json')
        .end(function(err, res) {
            var torrents = JSON.parse(res.text);
            this.list = torrents;

            this.trigger(this.list);
        }.bind(this));
    },
});

module.exports = Store;
