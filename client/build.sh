echo TIME TO BUILD
echo ...
rm -rf build
mkdir build
electron-packager . --all --out=build --overwrite

echo BUILD COMPLETE, UPLOADING TO ITCH
echo (keep an eye out, you may be prompted!)
echo ...
butler push build/ElderRune-win32-x64 ambushsabre/elder-runes:win64pre-alpha
butler push build/ElderRune-win32-ia32 ambushsabre/elder-runes:win32pre-alpha
butler push build/ElderRune-linux-ia32 ambushsabre/elder-runes:linux32pre-alpha
butler push build/ElderRune-linux-x64 ambushsabre/elder-runes:linux64pre-alpha

echo PUSH COMPLETE!
echo TIME TO PARTY!!!!!!!!!!
