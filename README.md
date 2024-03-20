# Simple Image Processor

Golang backend service with HTTP routes for image processing. This service has the following functionalities:

- Convert image files from PNG to JPEG.
- Resize images according to specified dimensions.
- Compress images to reduce file size while maintaining reasonable quality.

## Prerequisite

- This project is built using [**Go version 1.21.0**](https://go.dev/dl/), and it is expected to be developed using this specific version of Golang to ensure the desired outcome.
- Please ensure that OpenCV is installed on your computer, with a minimum version of [**OpenCV 4.7.0**](https://gocv.io/getting-started/).
 
## How to Run

- Install project dependencies using the command `go mod download`.
- Run the service using the command `go run .` or `go run main.go`.

## How to Build

Execute the following command to build the binary:
```shell
go build -ldflags "-s -w" -o ./out .
```

Then you can run the service using the command `./out`.

## Usage

### [POST] /converter
Convert image file from PNG to JPEG.

#### Request Body
|          Name | Required |  Type   | Description                                                                                                                                                           |
| -------------:|:--------:|:-------:| --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|     `image` | required | file  | The image file to be converted. <br/><br/> Supported MIME type: `image/png`.                                                                     |

#### Success response
The response will return an image file if successful.

#### Error response
The response will be returned as JSON in case of an error. For example:
```json
{
    "code": 400,
    "message": "You have to provide `image` field!",
    "status": "error"
}
```

### [POST] /resizer
Resize image according to specified dimensions.

#### Request Body
| Name | Required | Type | Description |
| ----:|:--------:|:----:| ----------- |
| `image` | required | file | The image file to be converted. <br/><br/> Supported MIME type: `image/png` and `image/jpeg`. |
| `width` | required | integer  | Width dimension. |
| `height` | required | integer  | Height dimension.|

#### Success response
The response will return an image file if successful.

#### Error response
The response will be returned as JSON in case of an error. For example:
```json
{
    "code": 400,
    "message": "You have to provide `image` field!",
    "status": "error"
}
```

### [POST] /compressor
Compress image to reduce file size while maintaining reasonable quality.

#### Request Body
| Name | Required | Type | Description |
| ----:|:--------:|:----:| ----------- |
| `image` | required | file | The image file to be converted. <br/><br/> Supported MIME type: `image/png` and `image/jpeg`. |

#### Success response
The response will return an image file if successful.

#### Error response
The response will be returned as JSON in case of an error. For example:
```json
{
    "code": 400,
    "message": "You have to provide `image` field!",
    "status": "error"
}
```

