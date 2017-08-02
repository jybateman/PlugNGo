PlugNGo Installation Guide
==========================
In this guide we will explain the steps to install both the servers of PlugNGo. PlugNGo consist of a server which will communicate with all SmartPlugs in proximity and serve as web server for the user, so that the user can quickly and easily control and monitor all SmartPlugs via a web interface.

Dependency
==========
PlugNGo depends on the following package:

 - Websocket: To communicate with the web interface in real time PlugNGo uses the experimental package `websocket golang.org/x/net/websocket`
 - MySQL: To store all SmartPlugs data in a database PlugNGo uses the MySQL to manage the database `github.com/go-sql-driver/mysql`
 - BLE (Soon to change): As of now we use a fork of paypal's GATT package to communicate with the SmartPLUG via BLE `github.com/currantlabs/gatt`, but this will surely be replaced by the BLE package `github.com/currantlabs/ble` when it is in a more stable state.

Server
======
The Server can be installed on any device that support BLE and must be in close enough proximity to be able to communicate with the SmartPLUG.
