Postfix bounce finder
=====================

Postfix logs parser to find and insert bounces into MongoDB.

* Compile

> go build bounce2db.go

* Run with default configuration file

> ./bounce2db /var/log/mail.info

* Run with other configuration file

> ./bounce2db -conf="/tmp/conf.json" /var/log/mail.info

* Run with many files

> ./bounce2db /var/log/mail.info /var/log/mail.info.1

> ./bounce2db /var/log/mail.info*

