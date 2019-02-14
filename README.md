# Realtime Chat Server using GraphQL Subscription and Redis


## frontend
```bash
npx create-nuxt-app frontend
cd frontend
yarn add @nuxtjs/apollo graphql-tag
```

add settings to `nuxt.config.js`
```javascript
  modules: ['@nuxtjs/apollo'],

  apollo: {
    clientConfigs: {
      httpEndpoint: 'http://localhost:8080',
      wsEndpoint: 'ws://localhost:8080',
      websocketsOnly: true
    }
  },
```

## References
- [Real-time Chat with GraphQL Subscriptions in Go](https://outcrawl.com/go-graphql-realtime-chat)
- [GoとRedisにおける簡単なチャットアプリケーション](https://medium.com/eureka-engineering/go-redis-application-28c8c793a652)
- [Redis の Pub/Sub を使って Node.js + WebSocket のスケールアウトを実現する方法](https://blog.dakatsuka.jp/2011/06/19/nodejs-redis-pubsub.html)
- [Apollo inside of NuxtJS](https://github.com/nuxt-community/apollo-module)
- [GraphQL と Nuxt.js でチャットを作る](https://www.aintek.xyz/posts/graphql-nuxt)