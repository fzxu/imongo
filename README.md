imongo
=======

An image server written in go lang and uses mongodb as backend persistence.

It makes good use of the 16M limitation of mongodb document in order to serve
lots of small files. (It does not use the GridFS because all the files are
smaller than 16M, using GridFS would be a waste of space)


You can send POST request to store a image and GET it back in any size.
The image would be store in mongodb based on `stored_size` from config.json.
