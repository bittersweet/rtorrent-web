import React from "react";

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
                            <div className="input-field">
                                <input id="search" type="search" required onChange={this.filterTorrents} name="search" placeholder="search" />
                                <label htmlFor="search"><i className="mdi-action-search"></i></label>
                                <i className="mdi-navigation-close"></i>
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

module.exports = Menu;
