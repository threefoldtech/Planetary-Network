#!/bin/bash
FILE=LICENSE

if [ ! -f "$FILE" ]; then
    echo "Please run script from main directory"
    exit 1
fi

rm -Rf src/moc*
rm -Rf src/vendor
rm -Rf *.deb
rm -Rf src/deploy
rm -Rf src/rcc*