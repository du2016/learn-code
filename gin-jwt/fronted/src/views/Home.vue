<template>
  <div class="home">
    <button @click="sayhello">sayhello</button><br>
    token: {{ token }}<br>
    msg: {{ msg }}<br>
    <button @click="changetoken">changetoken</button><br>
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from '@/components/HelloWorld.vue';

export default {
  name: 'home',
  data(){
    return {
        token: "",
        msg: ""
    }
  },
  methods:{
    sayhello: function(){
        var _this = this
        this.$http.get("/auth/hello",{headers: {'Authorization': 'Bearer ' + this.token,}}).then(function(response){
            _this.msg=response.data
        })
    },
    changetoken: function(){
        var _this= this
        this.$http.get("/auth/refresh_token",{headers: {'Authorization': 'Bearer ' + this.token,}}).then(function(response){
            localStorage.setItem("Authorization",response.data.token)
            _this.token = localStorage.getItem('Authorization');
        })
    }
  },
  mounted: function(){
    this.token = localStorage.getItem('Authorization');
  }
};
</script>
