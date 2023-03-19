# Build stage (AS stage name)
FROM golang:1.19.7-alpine3.17 AS builder
# define current image working directory 
WORKDIR /app
# copy all files from current folder to current image working dir
COPY . .
# build app to a single executable file
# -o: output; main: name of output binary file; main.go: main entry point file of the app
RUN go build -o main main.go


# Run stage (convert file to multistage (reduce the image size))
FROM alpine:3.17
WORKDIR /app
# copy the executable binary file from the builder stage to the run stage image
# /app/main: path to the file we want to copy; .: target location in the final image to copy that file to (i.e. /app)
COPY --from=builder /app/main .

# note down the port which is intended to be published
EXPOSE 8080
# define the default command to run when the container starts (run the executable file)
CMD [ "/app/mainn" ]