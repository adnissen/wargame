echo TIME TO BUILD
echo ...
rm -rf build
mkdir build
electron-packager . --all --out=build --overwrite

echo BUILD COMPLETE, UPLOADING TO ITCH
echo keep an eye out, you may be prompted!
echo ...
butler push build/ElderRunes-win32-ia32 ambushsabre/elder-runes:win-pre-alpha
butler push build/ElderRunes-linux-ia32 ambushsabre/elder-runes:linux-pre-alpha
butler push build/ElderRunes-darwin-x64 ambushsabre/elder-runes:mac-pre-alpha

echo PUSH COMPLETE!
echo TIME TO PARTY!!!!!!!!!!
