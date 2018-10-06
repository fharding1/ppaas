![Party Parrot](https://github.com/fharding1/ppaas/blob/master/parrot.gif?raw=true)![Smile Parrot](https://github.com/fharding1/ppaas/blob/master/smile_parrot.gif?raw=true)

# Party Parrot as a Service (PPaaS)

For all your serverless web scale party parrot needs.

https://cultofthepartyparrot.com/

## Setup

    dep ensure

## Usage

    go run main.go
    curl localhost:8080 -F "image=@/path/to/image.png" --output out.gif
