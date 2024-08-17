#!/bin/bash
python -u ./python_interface/main.py &
sleep 15
./s0backend
