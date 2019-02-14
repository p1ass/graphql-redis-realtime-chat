<template>
  <section class="container">
    <h1>ユーザーを作成する</h1>
    <div class="input-wrapper">
      <input 
        v-model="user" 
        type="text">
      <button @click="createUser">作成</button>
    </div>
  </section>
</template>

<script>
import MCreateUser from '@/apollo/mutations/createUser.gql'
import { mapActions } from 'vuex'

export default {
  data() {
    return {
      user: ''
    }
  },

  methods: {
    ...mapActions('users', ['storeUser']),

    async createUser() {
      if (this.user) {
        try {
          const res = await this.$apollo.mutate({
            mutation: MCreateUser,
            variables: {
              user: this.user
            }
          })

          // Success creating user
          if (res.data.createUser === this.user) {
            this.storeUser(this.user)
            this.$router.push('/room')
          }
        } catch (e) {
          console.log(e)
          alert('既にこのユーザー名は使用されています')
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/style/global.scss';

.container {
  margin: 16px;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.input-wrapper {
  margin-top: 16px;
  height: 40px;
}

input {
  height: 100%;
  margin-right: 16px;
  padding-left: 8px;
  border: 1px solid $border-color;
  border-radius: $border-radius;
  font-size: 20px;
}

button {
  width: 64px;
  height: 100%;
  color: white;
  background: $color-blue;
  border: 0;
  border-radius: $border-radius;
  font-size: 16px;
}
</style>
