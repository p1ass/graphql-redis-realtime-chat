export const state = () => ({
  user: ''
})

export const getters = {}

export const mutations = {
  setUser(state, user) {
    state.user = user
  }
}

export const actions = {
  storeUser({ commit }, user) {
    commit('setUser', user)
  }
}
