FROM node:18.8.0-alpine

RUN mkdir app

WORKDIR /app

RUN npm i -g http-server

COPY ./package*.json ./

RUN npm install

ARG VITE_API_URL \
    PORT=8080

ENV VITE_API_URL=${VITE_API_URL}  \
    PORT=${PORT}

COPY . .

RUN npm run build

EXPOSE 80

CMD ["http-server", "dist", "-p", "80", "-a", "0.0.0.0"]