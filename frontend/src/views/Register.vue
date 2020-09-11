<template>
  <div class="register">
    <img alt="Vue logo" src="../assets/logo.png">
   <!-- <HelloWorld msg="Welcome to Your Vue.js App"/> -->
   <h1>chatapp register page</h1>
   <p>{{ msg }}</p>
  </div>
</template>

<script>
// @ is an alias to /src
//import HelloWorld from '@/components/HelloWorld.vue'

export default {
  name: 'Register',
  components: {
    //HelloWorld
  },
  data:function(){
    return{
      msg:String
    }
  },
  created:function(){
    this.msg = "gggggg"
  }
}
</script>
