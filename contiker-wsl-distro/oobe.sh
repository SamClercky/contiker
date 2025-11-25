#!/bin/bash

read -p "Do you want to clone the contiki-ng repository? [yn] " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo "Cloning from https://github.com/contiki-ng/contiki-ng.git"

    git clone --recurse-submodules --shallow-submodules --depth 1 https://github.com/contiki-ng/contiki-ng.git ~/contiki-ng
fi
