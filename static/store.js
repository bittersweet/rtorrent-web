import Reflux from "reflux";
// https://github.com/visionmedia/superagent/wiki/Superagent-for-Webpack
var request = require("superagent");

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
        .get(`${ hostName }/torrents.json`)
        .end(function(err, res) {
            var torrents = JSON.parse(res.text);
            this.list = torrents;

            this.trigger(this.list);
        }.bind(this));
    },

    onChangeStatus: function(torrent, newStatus) {
        var hash = torrent.props.data.hash;
        var url = `${ hostName }/torrents/${ hash }/changestatus?status=${ newStatus }`;
        request.post(url, function(){});
    },

    onCopyFiles: function(torrent) {
        var hash = torrent.props.data.hash;
        var url = `${ hostName }/torrents/${ hash }/copy`;
        request.get(url, function(){});
    },
});

module.exports = Store;
