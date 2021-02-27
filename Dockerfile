FROM golang:1.16.0-buster as build-stage

WORKDIR /src
COPY src/ /src
RUN go build -o /app/weather_center


FROM debian:buster as product-stage
EXPOSE 8080

RUN groupadd -r wcenter
RUN useradd -r -u 1010 -g wcenter wcenter

WORKDIR /app
COPY --from=build-stage /app/weather_center /app/weather_center
RUN mkdir /app/db && \ 
	chown -R wcenter /app

USER wcenter
ENTRYPOINT ["/app/weather_center"]