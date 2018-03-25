FROM alpine

ARG TOKEN
ENV TELEGRAM_TOKEN $TOKEN
ARG PORT
ENV PORT $PORT

WORKDIR /app
COPY ./FoodHelperBot_unix /app/

CMD ./FoodHelperBot_unix