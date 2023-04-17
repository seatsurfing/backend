FROM node:18-alpine AS commons-builder
WORKDIR /app
RUN apk add git
RUN git init && \
    git remote add -f origin https://github.com/seatsurfing/backend.git && \
    git config core.sparseCheckout true && \
    echo 'commons' >> .git/info/sparse-checkout && \
    git pull origin master
WORKDIR /app/commons/ts
RUN npm install
RUN npm run build

FROM node:18-alpine AS builder
COPY --from=commons-builder /app/commons/ts/ /app/commons/ts
ADD . /app/
WORKDIR /app
RUN npm install
RUN npm install --save ./commons/ts
RUN npm run build

FROM node:18-alpine
ARG CI_VERSION
ENV NEXT_PUBLIC_PRODUCT_VERSION=$CI_VERSION
ENV NODE_ENV=production
ENV PORT=3001
WORKDIR /app
COPY --from=builder /app/public ./public
COPY --from=builder --chown=65532:65532 /app/build/standalone ./
COPY --from=builder --chown=65532:65532 /app/build/static ./build/static
RUN mkdir "/.npm" && chown -R 65532:65532 "/.npm"
EXPOSE 3001
USER 65532:65532
CMD ["node", "server.js"]