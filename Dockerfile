FROM alpine

ARG CONFIG
ARG PORT
ARG CERT
ARG KEY

WORKDIR /app
COPY ./FoodHelperBot_unix /app/

CMD ./FoodHelperBot_unix --config=$CONFIG --port=$PORT --cert=$CERT --key=$KEY
