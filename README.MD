fbAlbumBackup
======================
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/toomore/gogrs/master/LICENSE)

Facebook album tool that backup all photos in albums and supports concurrent download.

Install
--------------

    go get -u -x github.com/howeih/fbAlbumBackup

Note, you need to go to [faceook developer page](https://developers.facebook.com/tools/explorer?method=GET&path=me) and get token.

Usage
---------------------

    fbAlbumBackup [options] 

Options
---------------

- `-token`  
- `-worker` number of workers. Default workers is 5


Examples
---------------

Backup all photos with 10 workers.

  fbAlbumBackup -token=access-token -worker=10
