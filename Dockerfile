FROM alpine

ARG TOKEN
ARG YELP
ARG PORT
ARG CERT
ARG KEY

WORKDIR /app
COPY ./FoodHelperBot_unix /app/

CMD ./FoodHelperBot_unix --token=$TOKEN --yelpKey=$YELP --port=$PORT --cert=$CERT --key=$KEY