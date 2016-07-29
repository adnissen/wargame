echo TIME TO BUILD
rm -rf build
mkdir build
electron-packager . --all --out=build --overwrite
