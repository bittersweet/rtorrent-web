import React from "react";
import Menu from "./menu";
import Statistics from "./statistics";
import Torrent from "./torrent";
import TorrentList from "./torrentlist";

var esHost = "http://localhost:8001";
var hostName = "http://localhost:8000";

// http://stackoverflow.com/questions/10420352/converting-file-size-in-bytes-to-human-readable
Object.defineProperty(Number.prototype,'fileSize',{value:function(a,b,c,d){
 return (a=a?[1e3,'k','B']:[1024,'K','iB'],b=Math,c=b.log,
 d=c(this)/c(a[0])|0,this/b.pow(a[0],d)).toFixed(2)
 +' '+(d?(a[1]+'MGTPEZY')[--d]+a[2]:'Bytes');
},writable:false,enumerable:false});

var App = React.createClass({
    filter: function(filter) {
        this.setState({filterOn: filter});
    },

    searchTorrents: function(query) {
        this.setState({queryOn: query});
    },

    getInitialState: function() {
        return {
            filterOn: "none",
            queryOn: "",
        };
    },

    render: function() {
        return (
            <div>
                <Menu filter={this.filter} search={this.searchTorrents} currentFilter={this.state.filterOn} />
                <TorrentList hostName={hostName} pollInterval={2000} filterOn={this.state.filterOn} queryOn={this.state.queryOn} />
            </div>
        )
    }
});

var client = new EventSource(esHost);
client.onmessage = function(message) {
    Materialize.toast(message.data, 3000);
}

React.render(
    <App />,
    document.getElementById('content')
);

