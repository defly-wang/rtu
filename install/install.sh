#!/bin/sh

#make dome dir
mkdir /etc/rtu
mkdir /opt/rtu

cp bin/* /opt/rtu/
chmod a+x /opt/rtu/*

cp etc/* /etc/rtu/ -rf

cp cgi/* /usr/boa/www/cgi-bin/
chmod a+x /usr/boa/www/cgi-bin/rtucgi.cgi


cp autorun/rc.S /etc/init.d/

reboot
