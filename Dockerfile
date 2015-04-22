# This image includes multiple ONBUILD triggers which should cover most applications. The 
# build will COPY . /usr/src/app, RUN go get -d -v, and RUN go install -v.

# This image also includes the CMD ["app"] instruction which is the default command when 
# running the image without arguments.
FROM golang:1.4.2-onbuild