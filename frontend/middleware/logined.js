export default function({ store, redirect }) {
  // ユーザーが認証されていないとき
  if (!store.state.users.user) {
    return redirect('/')
  }
}
