set -e
rm -rf demo
mkdir -p demo
cd bugs/vnet/
echo "Building vnet"
go build -o ../../demo/vnet_demo
cd ../main
echo "Building L8Bugs services"
go build -o ../../demo/bugs_demo
cd ../website/main1
echo "Building ui"
go build -o ../../../demo/ui_demo
cd ..
cp -r ./web ../../demo/.
cd ../../demo
./vnet_demo &
sleep 1
./bugs_demo &
./ui_demo

pkill demo
cd ..
rm -rf demo
