FROM alpine

ARG TOKEN
ARG PORT
ARG CERT
ARG KEY

WORKDIR /app
COPY ./FoodHelperBot_unix /app/

CMD ./FoodHelperBot_unix --token=$TOKEN --port=$PORT --cert=$CERT --key=$KEY