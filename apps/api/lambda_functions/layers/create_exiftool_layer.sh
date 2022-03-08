#!/bin/sh

export EXIFTOOL_VERSION=12.40
export PERL_VERSION="5-30"

rm -rf layer
curl -sS https://shogo82148-lambda-perl-runtime-ap-southeast-2.s3.amazonaws.com/perl-${PERL_VERSION}-runtime.zip > perl.zip
mkdir layer
cd layer
unzip ../perl.zip
cd ..
rm perl.zip
curl https://www.sno.phy.queensu.ca/~phil/exiftool/Image-ExifTool-${EXIFTOOL_VERSION}.tar.gz | tar -xJ
mkdir -p layer/bin
cp Image-ExifTool-${EXIFTOOL_VERSION}/exiftool layer/bin/.
sed -i "" "1 s/^.*$/#\!\/opt\/bin\/perl -w/" layer/bin/exiftool
cp -r Image-ExifTool-${EXIFTOOL_VERSION}/lib layer/bin/.
rm -rf Image-ExifTool-${EXIFTOOL_VERSION}
cd layer
zip -r layer.zip ./*
cd ..
mv layer/layer.zip .
rm -rf layer
