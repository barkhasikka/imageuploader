# imageuploader

Create a service for storing and retrieving image data. The service should follow best practices of a web service, examples:

* data validation
* correct status codes
* proper headers (like `Content-Type`, etc)
* ...

Feel free to use whatever database backend you like. For programming language we would like Go to be used.

The service should be **performant** and ideally **easily scalable**. To limit the scope of the work it's OK to make less than ideal choices as long as they are motivated. For instance "I used the Windows Registry on my computer as database, it won't scale beyond three entries but if I wanted to make it scalable I would change to `<insert buzzword here>`".

We want the code with a nice commit history as a git bundle (``git bundle create <bundle.name> --all``) and if you can deploy it somewhere that would be great but it's not a requirement. Any questions can be sent to gustaf@apple.com.

## Endpoints

`GET /v1/images`

List metadata for stored images.

`GET /v1/images/<id>`

Get metadata for image with id `<id>`.

`GET /v1/images/<id>/data` 


Get image data for image with id `<id>`.

Optional GET parameter: `?bbox=<x>,<y>,<w>,<h>` to get a cutout of the image.

`POST /v1/images`

Upload new image. Request body should be image data.

`PUT /v1/images/<id>`

Update image. Request body should be image data.

## Image metadata

Image metadata should be the following fields:

* Filesize of the image
* Image dimensions
* Image type (gif, jpg etc)
* Date/time image was uploaded
* Extra fields of your choice

It is up to you how you want to represent the image metadata in your API.

Some notes about the assignment:

To test update:
curl -X PUT -F 'file=@/path/to/image/file' localhost:8089/v1/images/:id -H "Content-Type: multipart/form-data"

In ideal scenarios we should convert all the images into one of png/ jpg to make code cleaner and handling of images much easier.
Many tools like ffmpeg can be a good choice to convert every image to certain types depending on the requirement of course.