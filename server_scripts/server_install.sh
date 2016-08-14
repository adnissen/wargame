cp .profile ~/.profile
cd ~
curl -O https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz
tar xvf go1.6.linux-amd64.tar.gz
sudo chown -R root:root ./go
sudo mv go /usr/local

source ~/.profile
