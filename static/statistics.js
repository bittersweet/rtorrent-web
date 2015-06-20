import React from "react";

var Statistics = React.createClass({
    updateCounts: function() {
        var up = 0;
        var down = 0;
        this.props.torrents.map(function(torrent) {
            up += torrent.get_up_rate_raw;
            down += torrent.get_down_rate_raw;
        });
        this.setState({up: up.fileSize(), down: down.fileSize()})
    },

    getInitialState: function() {
        return {
            up: 0, down: 0
        };
    },

    componentWillReceiveProps: function() {
        this.updateCounts()
    },

    render: function() {
        return (
            <div className="statistics">upload: {this.state.up}, download: {this.state.down}</div>
        )
    },
});

module.exports = Statistics;
