# cvdhelper
Computer Vision Data Helper (cvdhelper) is a set of functions I find that I use a lot in analyzing data before pre-processing it.

## Disclaimer ---
If you are looking for a lot of image tools there are some way better ones out there.  
This is not a super extensive set of functions.  Batching a slice of image.Image into a 4d tensor of type float32.   Mainly just turning image.Images to NCHW or NHWC tensors. Finding the average, max, min values.  Some basic add, divide, multiply to all elements. Cloning, Mirroring tensors.  You know stuff that can be useful when putting images to train a cnn.  Tensor to image.Image. 

Note: Even though this package doesn't have resizing. (easly found in other image packages).  If a slice of images of different sizes is used to make a 4d tensor. It will pad the smaller h, w dims with zeros. You can see in the outputfromtesting (re imaged from a 4d tensor) folder and testimgs (original images) folder.