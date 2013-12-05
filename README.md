imongo
=======

An image server written in go lang and uses mongodb as backend persistence.

It makes good use of the 16M limitation of mongodb document in order to serve
lots of small files. (It does not use the GridFS because all the files are
smaller than 16M, using GridFS would be a waste of space)


You can send POST request to store a image and GET it back in any size.
The image would be store in mongodb based on `stored_size` from config.json.

You can post an image to a path(case insensitive) with the image binary data, e.g.

`http://localhost/foo/bar/picture.jpg`

Then you can get the image from different sizes:

`http://localhost/foo/bar/picture.jpg` will get the purely stored image(The max size is the stored_size from configuration).

`http://localhost/foo/bar/picture__100x200.jpg` will resize the image to 100 width and 200 height. This could change the aspect ratio.

`http://localhost/foo/bar/picture__0x200.jpg` will resize the image to 200 height and keep the aspect ratio.

`http://localhost/foo/bar/picture__100x0.jpg` will resize the image to 100 width and keep the aspect ratio.

`http://localhost/foo/bar/picture__100z200.jpg` will resize and crop to 100 width and 200 height.

All the generated images would be stored in the cache folder(configurable). 
You can use any proxy(e.g. nginx) to serve the data next time.