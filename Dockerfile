FROM golang:1.19.7-alpine3.17

# define current image working directory 
WORKDIR /app

# copy all files from current folder to current image working dir
COPY . .

# build app to a single executable file
# -o: output; main: name of output binary file; main.go: main entry point file of the app
RUN go build -o main main.go


# note down the port which is intended to be published
EXPOSE 8080
# define the default command to run when the container starts (run the executable file)
CMD [ "/app/mainn" ]