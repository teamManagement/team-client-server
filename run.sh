#!/bin/bash
echo "stop teamClientServer..."
$1 -cmd=stop

echo "uninstall teamClientServer..."
$1 -cmd=uninstall

echo "install teamClientServer..."
$1 -cmd=install

if [ $? -ne 0 ]; then
    exit $?
fi

echo "start teamClientServer"
$1 -cmd=start
