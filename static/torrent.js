import React from "react";

import Actions from "./actions";

var Torrent = React.createClass({
    onClick: function() {
        var url = hostName + '/torrents/' + this.props.data.hash;
        window.open(url, '_blank');
    },

    changeStatus: function(event) {
        event.preventDefault();
        Actions.changeStatus(this, event.target.text)
    },

    copyFiles: function(event) {
        event.preventDefault();
        Actions.copyFiles(this);
    },

    render: function() {
        var up_total = this.props.data.get_up_total.fileSize()
        var state;
        switch(this.props.data.state) {
            case 0:
                state = "stopped"
                break;
            case 1:
                state = "started"
                break;
        }

        return (
            <tr className="torrent">
                <td>{this.props.data.tracker}</td>
                <td>{state}</td>
                <td>
                <span onClick={this.onClick}>{this.props.data.name}</span>
                { ' ' }
                <a href="#" onClick={this.changeStatus}>stop</a>
                { ' ' }
                <a href="#" onClick={this.changeStatus}>start</a>
                { ' ' }
                <a href="#" onClick={this.changeStatus}>remove</a>
                { ' ' }
                <a href="#" onClick={this.copyFiles}>copy</a>
                </td>
                <td className="center">{this.props.data.size_files}</td>
                <td className="done_total center">
                  {this.props.data.bytes_done} / {this.props.data.size_bytes} ({this.props.data.percentage_done})
                </td>
                <td className="center">{this.props.data.get_down_rate}</td>
                <td className="center">{this.props.data.get_up_rate}</td>
                <td className="center">{up_total}</td>
                <td className="center">{this.props.data.ratio}</td>
                <td className="center">{this.props.data.peers_connected}</td>
            </tr>
        );
    }
});

module.exports = Torrent;
