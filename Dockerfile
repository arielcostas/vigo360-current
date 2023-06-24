FROM golang:1.20.5 AS build

# prebuild() {
# 	mkdir -p assets/{extra,images,papers,profile,thumb}
# 	cp static/* assets/
# 	sass --source-map -s compressed styles/:assets/
# 	chmod -R a+r assets/
# }
# 
# if [ "$1" == "run" ];
# then
# 	prebuild
# 	export $(cat .env | grep -v '^#' | xargs)
# 	go run -ldflags "-X main.version=$(git rev-parse --short HEAD)" .
# elif  [ "$1" == "build" ];
# then
# 	prebuild
# 	go build -o vigo360 -ldflags "-X main.version=$(git rev-parse --short HEAD)" .
# else
# 	printf "Vigo360 launcher script\nusage: ${0} [command]\n\n	build: compiles vigo360 to a binary to be ran in production\n	  run: runs vigo360 via in memory.\n"
# fi
# Copilot, generate a Dockerfile that does this

COPY . /app
WORKDIR /app
RUN go build -o vigo360-server .

FROM gcr.io/distroless/base-debian10 AS runtime

COPY --from=build /app/vigo360-server /app/executable

WORKDIR /app

EXPOSE 8080

CMD ["/app/executable"]


