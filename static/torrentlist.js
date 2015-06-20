import React from "react";
import Statistics from "./statistics";
import Torrent from "./torrent";

var TorrentList = React.createClass({
    loadTorrents: function() {
        $.ajax({
            url: this.props.hostName + '/torrents.json',
            dataType: 'json',
            success: function(data) {
                this.setState({data: data});
            }.bind(this)
        });
    },

    getInitialState: function() {
        return {
            data: [],
            sortOn: "default",
        };
    },

    componentDidMount: function() {
        this.loadTorrents();
        setInterval(this.loadTorrents, this.props.pollInterval);
    },

    sort: function(object) {
        sortDirection = object.target.id;
        if (this.state.sortOn != sortDirection) {
            this.setState({sortOn: sortDirection});
        } else {
            this.setState({sortOn: "default"});
        }
    },

    render: function() {
        var filterOn = this.props.filterOn;
        var sortOn = this.state.sortOn;
        var queryOn = this.props.queryOn;

        var torrents = this.state.data.map(function(torrent) {
            return (
                <Torrent key={torrent.hash} data={torrent} />
            );
        });

        if (queryOn != "") {
            torrents = torrents.filter(function(torrent) {
                torrent = torrent.props.data;
                return (torrent.name.toLowerCase().indexOf(queryOn) >= 0)
            });
        }

        torrents = torrents.filter(function(torrent) {
            torrent = torrent.props.data;
            switch(filterOn) {
                case "uploads":
                if (parseInt(torrent.get_up_rate) > 0) {
                    return (
                        <Torrent key={torrent.hash} data={torrent} />
                    );
                }
                break;
                case "downloads":
                if (parseInt(torrent.get_down_rate) > 0) {
                    return (
                        <Torrent key={torrent.hash} data={torrent} />
                    );
                }
                break;
                case "none":
                return (
                    <Torrent key={torrent.hash} data={torrent} />
                );
                break;
            }
        });

        if (sortOn != "default") {
            torrents = torrents.sort(function(a, b) {
                switch(sortOn) {
                    case "name":
                        var nameA = a.props.data.name.toLowerCase();
                        var nameB = b.props.data.name.toLowerCase();
                        if (nameA < nameB) {
                            return -1;
                        }
                        return 0;
                        break;
                    case "download":
                        return b.props.data.get_down_rate_raw - a.props.data.get_down_rate_raw;
                        break;
                    case "upload":
                        return b.props.data.get_up_rate_raw - a.props.data.get_up_rate_raw;
                        break;
                    case "up_total":
                        return b.props.data.get_up_total - a.props.data.get_up_total;
                        break;
                    case "ratio":
                        return b.props.data.ratio - a.props.data.ratio;
                        break;
                }
            });
        }

        return (
            <div>
                <Statistics data={this.state.data} />
                <table className="torrentList striped">
                    <thead>
                    <tr>
                        <th>Tracker</th>
                        <th>Status</th>
                        <th id="name" onClick={this.sort}>Name</th>
                        <th className="center">Files</th>
                        <th className="done_total center">Done / Total</th>
                        <th id="download" className="center" onClick={this.sort}>Down rate</th>
                        <th id="upload" className="center" onClick={this.sort}>Up rate</th>
                        <th id="up_total" className="center" onClick={this.sort}>Up total</th>
                        <th id="ratio" className="center" onClick={this.sort}>Ratio</th>
                        <th className="center">Peers Connected</th>
                    </tr>
                    </thead>
                    <tbody>
                    {torrents}
                    </tbody>
                </table>
            </div>
        );
    }
});

module.exports = TorrentList;
