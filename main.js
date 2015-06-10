var app = require('app');
var BrowserWindow = require('browser-window');

require('crash-reporter').start();

var mainWindow = null;

app.on('window-all-closed', function() {
    if (process.platform != 'darwin') {
        app.quit();
    }
});

app.on('ready', function() {
    var options = {
        "width": 1920,
        "height": 1024,
        "node-integration": false // so $ works, https://github.com/atom/electron/issues/254
    }
    mainWindow = new BrowserWindow(options);
    mainWindow.loadUrl('file://' + __dirname + '/static/index.html');
    mainWindow.openDevTools();

    mainWindow.on('closed', function() {
        mainWindow = null;
    });
});
