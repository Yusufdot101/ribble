######################
# dependencies layer #
######################
FROM node:24.15-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci


###############
# build layer #
###############
FROM node:24.15-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG NEXT_PUBLIC_BASE_USER_SERVICE_API_URL
ARG NEXT_PUBLIC_BASE_CHAT_SERVICE_API_URL

ENV NEXT_PUBLIC_BASE_USER_SERVICE_API_URL=$NEXT_PUBLIC_BASE_USER_SERVICE_API_URL
ENV NEXT_PUBLIC_BASE_CHAT_SERVICE_API_URL=$NEXT_PUBLIC_BASE_CHAT_SERVICE_API_URL


RUN npm run build


#################
# runtime layer #
#################
FROM node:24.15-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production

# Only copy what's needed to run
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

EXPOSE 3000
USER node
CMD ["node", "server.js"]
