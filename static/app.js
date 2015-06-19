// http://stackoverflow.com/questions/10420352/converting-file-size-in-bytes-to-human-readable
Object.defineProperty(Number.prototype,'fileSize',{value:function(a,b,c,d){
 return (a=a?[1e3,'k','B']:[1024,'K','iB'],b=Math,c=b.log,
 d=c(this)/c(a[0])|0,this/b.pow(a[0],d)).toFixed(2)
 +' '+(d?(a[1]+'MGTPEZY')[--d]+a[2]:'Bytes');
},writable:false,enumerable:false});

// var hostname = 'http://192.168.2.7:8000';
var hostname = 'http://localhost:8000';

var Menu = React.createClass({
    filterUploads: function() {
        this.props.filter('uploads');
    },

    filterDownloads: function() {
        this.props.filter('downloads');
    },

    removeFilters: function() {
        this.props.filter('none');
    },

    filterTorrents: function(event) {
        this.props.search(event.target.value);
    },

    getClassName: function(item) {
        var cF = this.props.currentFilter;
        if (item == "uploading" && cF == "uploads") {
                return "active";
        }
        if (item == "downloading" && cF == "downloads") {
            return "active";
        }
        if (item == "all" && cF == "none") {
            return "active";
        }
    },

    render: function() {
        return (
            <nav>
                <div className="nav-wrapper menu">
                    <a href="#" className="brand-logo">Rtorrent-Web</a>
                    <ul id="nav-mobile" className="right hide-on-med-and-down">
                        <li>
                            <div className="search">
                                <input id="search" onChange={this.filterTorrents} name="search" placeholder="search" type="text" />
                            </div>
                        </li>
                        <li className={this.getClassName('uploading')}>
                            <a href="#" onClick={this.filterUploads}>Uploading only</a>
                        </li>
                        <li className={this.getClassName('downloading')}>
                            <a href="#" onClick={this.filterDownloads}>Downloading only</a>
                        </li>
                        <li className={this.getClassName('all')}>
                            <a href="#" onClick={this.removeFilters}>Show all</a>
                        </li>
                    </ul>
                </div>
            </nav>
        )
    }
});

var Statistics = React.createClass({
    updateCounts: function() {
        var up = 0;
        var down = 0;
        this.props.data.map(function(torrent) {
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
            queryOn: ""
        };
    },

    render: function() {
        return (
            <div>
                <Menu filter={this.filter} search={this.searchTorrents} currentFilter={this.state.filterOn} />
                <TorrentList pollInterval={2000} filterOn={this.state.filterOn} queryOn={this.state.queryOn} />
            </div>
        )
    }
});

var TorrentList = React.createClass({
    loadTorrents: function() {
        $.ajax({
            url: hostname + '/torrents.json',
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

var Torrent = React.createClass({
    onClick: function() {
        var url = hostname + '/torrents/' + this.props.data.hash;
        window.open(url, '_blank');
    },

    changeStatus: function(event) {
        event.preventDefault();

        var url = hostname + '/torrents/' + this.props.data.hash + '/changestatus?status=' + event.target.text;

        $.get(url, function(data) {
            // console.log(data);
        });
    },

    copyFiles: function(event) {
        event.preventDefault();

        var url = hostname + '/torrents/' + this.props.data.hash + '/copy';

        $.get(url, function(data) {
            console.log(data);
        });
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

var client = new EventSource("http://192.168.2.7:8001");
client.onmessage = function(message) {
    Materialize.toast(message.data, 3000);
}

React.render(
    <App />,
    document.getElementById('content')
);

