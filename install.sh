# Alsa
sudo apt-get install libasound2 libasound2-dev -y
sudo pip install pyalsaaudio -y

# Aubio?
sudo apt-get install  python-numpy python-scipy python-matplotlib ipython ipython-notebook python-pandas python-sympy python-nose -y
sudo apt-get install build-essential python-dev git scons swig -y
#sudo apt-get install aubio-tools libaubio-dev libaubio-doc -y
#csudo pip install aubio

#WS2811
git clone https://github.com/jgarff/rpi_ws281x.git libs/neopixel-py
cd libs/neopixel-py
scons
sudo cp ws2811.h /usr/local/include/
sudo cp pwm.h /usr/local/include/
sudo cp rpihw.h /usr/local/include/
sudo cp libws2811.a /usr/local/lib/
cd python
sudo python setup.py install
cd ../../../

# Gvm
sudo apt-get install curl git mercurial make binutils bison gcc build-essential -y
curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer | zsh
source /home/pi/.gvm/scripts/gvm
gvm install go1.4
echo Use Go 1.4
gvm use go1.4
export GOROOT_BOOTSTRAP=$GOROOT
gvm install go1.9
gvm use go1.9 --default
go get github.com/cocoonlife/goalsa
go get github.com/simonassank/aubio-go
go get github.com/simonassank/go_ws2811
