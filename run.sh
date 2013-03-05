export PYTHONPATH=$PYTHONPATH:./CarDraw/src
(cd capture && make)
./capture/capture | python ./CarDraw/src/main/Main.py
