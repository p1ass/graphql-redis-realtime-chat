<template>
  <section class="container">

    <div class="dispaly">
      <div class="users-wrapper">
        <h2>ユーザ一覧</h2>
        <p 
          v-for="(user, index) in users" 
          :key="index"
          class="user-box" >{{ user }}</p>
      </div>

      <div class="messages-wrapper">
        <h2>メッセージ</h2>
        <p 
          v-for="(message, index) in messages" 
          :key="index"
          class="message-box">{{ message.user }} : {{ message.message }}</p>
      </div>
    </div>
		
    <div class="input-wrapper">
      <input
        v-model="message"
        type="text"
        placeholder="メッセージを入力する" 
        @keyup.enter="postMessage">
      <button @click="postMessage">送信</button>
    </div>
  </section>
</template>
	
<script>
import QUsers from '@/apollo/queries/users.gql'
import MPostMessage from '@/apollo/mutations/postMessage.gql'
import SUserJoined from '@/apollo/subscriptions/userJoined.gql'
import SMessagePosted from '@/apollo/subscriptions/messagePosted.gql'
import { mapState } from 'vuex'

export default {
  middleware: 'logined',

  data() {
    return {
      message: '',
      messages: []
    }
  },

  computed: {
    ...mapState('users', ['user'])
  },

  methods: {
    async postMessage() {
      const res = await this.$apollo.mutate({
        mutation: MPostMessage,
        variables: {
          user: this.user,
          message: this.message
        }
      })

      console.log(res)
    }
  },
  apollo: {
    users: {
      query: QUsers,
      subscribeToMore: {
        document: SUserJoined,
        variables() {
          return {
            user: this.user
          }
        },

        updateQuery: (prev, { subscriptionData }) => {
          // Response data not found
          if (!subscriptionData.data) {
            return prev
          }

          const user = subscriptionData.data.userJoined
          if (prev.users.find(u => u === user)) {
            return prev
          }

          return Object.assign({}, prev, {
            users: [user, ...prev.users]
          })
        }
      }
    },

    // Use Simple subscription becase messagePosted doesn't involve with any query.
    $subscribe: {
      messagePosted: {
        query: SMessagePosted,

        variables() {
          return {
            user: this.user
          }
        },

        result(res) {
          // save data to store
          this.messages.unshift(res.data.messagePosted)
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/style/global.scss';
.container {
  margin: 16px 0;
  color: rgb(50, 50, 50);
  height: 100vh - calc($margin * 2);
}

.dispaly {
  display: flex;
  height: 90vh;
  overflow: hidden;

  .users-wrapper {
    width: 20vw;
    border-right: 1px solid $border-color;
    justify-content: center;
    .user-box {
      margin: 16px;
      padding: 4px;
      height: 16px;
    }
  }

  .messages-wrapper {
    margin: 0 16px;
    flex: 1;
    overflow: hidden;

    .message-box {
      margin: 16px;
      padding: 8px;
      word-wrap: break-word;

      border: 1px solid $border-color;
      border-radius: $border-radius;
    }
  }
}

h2 {
  text-align: center;
  margin: 0 auto;
}

.input-wrapper {
  height: 48px;

  position: fixed;
  left: 16px;
  right: 16px;
  bottom: 16px;
  display: flex;
  justify-content: center;
  text-align: center;

  input {
    flex: 1;
    margin-right: 16px;
    padding-left: 8px;
    border: 1px solid $border-color;
    border-radius: $border-radius;
    font-size: 20px;
  }

  button {
    width: 64px;
    color: white;
    background: $color-blue;
    border: 0;
    border-radius: $border-radius;
    font-size: 16px;
  }
}
</style>
