<template>
  <div class="chat">
    <img alt="Vue logo" src="../assets/logo.png">
   <!-- <HelloWorld msg="Welcome to Your Vue.js App"/> -->
   <button v-on:click="openws"> open connection</button>
   <h1>chatapp chat page</h1>
   <p>{{ msg }}</p>
   <p>{{ msg2 }}</p>
   <p>{{ loc }}</p>
  </div>
</template>

<script>
// @ is an alias to /src
//import HelloWorld from '@/components/HelloWorld.vue'

export default {
  name: 'chat',
  components: {
    //HelloWorld
  },
  data:function(){
    return {
        msg:"",
        msg2:"",
        connection:null,
        loc : window.location
      }
  },
  methods:{
    openws: function(){
      this.connection = new WebSocket("ws://localhost:8881/echo")
      // Connection opened
      this.connection.onopen=function(){
        console.log("connection established")
        alert("connection established")
        this.connection.send("dddddd")
      }
      this.connection.onmessage=function(e){
        this.msg2=e.data
        this.msg=e
        alert(e)
      }
    }
  }
}
</script>
