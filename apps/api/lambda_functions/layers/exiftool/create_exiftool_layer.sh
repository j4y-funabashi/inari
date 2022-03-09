#!/bin/sh

export EXIFTOOL_VERSION=12.40
rm -rf layer
curl https://exiftool.org/Image-ExifTool-${EXIFTOOL_VERSION}.tar.gz  --output Image-ExifTool-${EXIFTOOL_VERSION}.tar.gz
tar -xf Image-ExifTool-${EXIFTOOL_VERSION}.tar.gz
rm -rf Image-ExifTool-${EXIFTOOL_VERSION}.tar.gz
mkdir -p layer/bin
cp Image-ExifTool-${EXIFTOOL_VERSION}/exiftool layer/bin/.
sed -i "1 s/^.*$/#\!\/opt\/bin\/perl -w/" ./layer/bin/exiftool
chmod a+x ./layer/bin/exiftool
cp -r Image-ExifTool-${EXIFTOOL_VERSION}/lib layer/bin/.
rm -rf Image-ExifTool-${EXIFTOOL_VERSION}
