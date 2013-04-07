export PYTHONPATH=$PYTHONPATH:./CarDraw/src
(cd capture && make)
python ./CarDraw/src/main/Main.py
