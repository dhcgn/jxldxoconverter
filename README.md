[![Go](https://github.com/dhcgn/jxldxoconverter/actions/workflows/go.yml/badge.svg)](https://github.com/dhcgn/jxldxoconverter/actions/workflows/go.yml)

# jxl dxo converter

> Don't use yet, this is a dirty prototyp! 

## Why?

JXL is an excellent format for pictures, but DXO Optics Pro provides only JPEG, DNG and TIFF. This applications brings JXL support to DXO Optics Pro.

## Workflow

Image format wich can be encoded directly to JXL will be passed to the JXL executablefrom https://github.com/libjxl/libjxl, other format will be converted to png with https://imagemagick.org and then passed to the JXL executable.

You should select TIF 16 Bit as export format to get the most of JXL. If you export your images in JPGs, this application will do a **lossless JPG transcoding** with ~20% size reduction.


## How to

![image](https://user-images.githubusercontent.com/6566207/153045705-446cd3f5-40c9-4802-aef9-7170bae1a7ba.png)

![image](https://user-images.githubusercontent.com/6566207/153044941-d6fab923-d8a8-4aa5-a2ca-54646bc4028d.png)

![image](https://user-images.githubusercontent.com/6566207/153045078-3d22ca06-61e7-4f7c-998c-2f8b96a731fc.png)

![image](https://user-images.githubusercontent.com/6566207/153045498-b9fbc687-1458-442f-800a-e97172fef3b7.png)
