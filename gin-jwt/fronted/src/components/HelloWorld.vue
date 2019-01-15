<template>
  <div class="hello">
    <input v-model="name"><br>
    <input v-model="password"><br>
    <button @click="submit">登录</button>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  props: {
    msg: String,
  },
  data () {
    return {
      name: "admin",
      password: "admin"
    }
  },
  methods: {
    submit: function(){
      var parmas={"username": this.name,"password":this.password}
      var _this = this
      this.$http.post("/login",parmas).then(function(response){
        console.log(response.data.token)
        console.log(response.data.expire)
        localStorage.setItem("Authorization",response.data.token)
        _this.$router.push('/')
      })
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
