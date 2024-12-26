#!/bin/sh

sleep 1

if [ -e "/opt/rtu/rtuclient.new" ]; then
    rm /opt/rtu/rtuclient
    mv /opt/rtu/rtuclient.new /opt/rtu/rtuclient
    chmod a+x /opt/rtu/rtuclient
fi

if [ -e "/opt/rtu/iotread.new" ]; then
    rm /opt/rtu/iotread
    mv /opt/rtu/iotread.new /opt/rtu/iotread
    chmod a+x /opt/rtu/iotread
fi


if [ -e "/usr/boa/www/cgi-bin/rtucgi.cgi.new" ]; then
    rm /usr/boa/www/cgi-bin/rtucgi.cgi
    mv /usr/boa/www/cgi-bin/rtucgi.cgi.new /usr/boa/www/cgi-bin/rtucgi.cgi
    chmod a+x /usr/boa/www/cgi-bin/rtucgi.cgi
fi

/opt/rtu/iotread &
/opt/rtu/rtuclient &

