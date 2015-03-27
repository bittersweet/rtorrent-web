var TorrentList = React.createClass({
    loadTorrents: function() {
        $.ajax({
        url: 'http://localhost:8000/torrents',
        dataType: 'json',
        success: function(data) {
            this.setState({data: data});
        }.bind(this)
        });
    },

    getInitialState: function() {
        return {data: []};
    },

    componentDidMount: function() {
        this.loadTorrents();
        setInterval(this.loadTorrents, this.props.pollInterval);
    },

    render: function() {
        var torrents = this.state.data.map(function (torrent) {
            return (
                <Torrent data={torrent} />
            );
        });
        return (
            <table class="pure-table" className="torrentList pure-table pure-table-striped">
                <thead>
                <tr>
                    <th>Tracker</th>
                    <th>Name</th>
                    <th className="center">Files</th>
                    <th className="done_total center">Done / Total</th>
                    <th className="center">Down rate</th>
                    <th className="center">Up rate</th>
                    <th className="center">Up total</th>
                    <th className="center">Ratio</th>
                    <th className="center">Peers Connected</th>
                </tr>
                </thead>
                <tbody>
                {torrents}
                </tbody>
            </table>
        );
    }
});

var Torrent = React.createClass({
    onClick: function() {
        var url = 'http://localhost:8000/torrents/' + this.props.data.hash;
        window.open(url, '_blank');
    },

    render: function() {
        return (
            <tr className="torrent">
                <td>{this.props.data.tracker}</td>
                <td onClick={this.onClick}>{this.props.data.name}</td>
                <td className="center">{this.props.data.size_files}</td>
                <td className="done_total center">
                {this.props.data.bytes_done} / {this.props.data.size_bytes} ({this.props.data.percentage_done})
                </td>
                <td className="center">{this.props.data.get_down_rate}</td>
                <td className="center">{this.props.data.get_up_rate}</td>
                <td className="center">{this.props.data.get_up_total}</td>
                <td className="center">{this.props.data.ratio}</td>
                <td className="center">{this.props.data.peers_connected}</td>
            </tr>
        );
    }
});

React.render(
    <TorrentList pollInterval={2000} />,
    document.getElementById('content')
);

