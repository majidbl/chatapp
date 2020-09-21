<template>
  <div class="container">
    <div v-if="logged == false">
      <b-form class="mt-5" inline>
        <label class="sr-only" for="inline-form-input-username">Username</label>
        <b-input id="inline-form-input-username" placeholder="Username" v-model="username"></b-input>
          
          <b-button variant="primary" v-on:click="loggIn()">Loggin</b-button>
      </b-form>
    </div>
    <div v-else>
    <div>
      <b-button v-b-toggle.sidebar-footer>List of Contact</b-button>
      <b-sidebar id="sidebar-footer" aria-label="Sidebar with custom footer" no-header shadow>
        <template v-slot:footer="{ hide }">
          <div class="d-flex bg-dark text-light align-items-center px-3 py-2">
            <form>
              <div class="form-group">
                <label for="username">UserName</label>
                <input
                  type="text"
                  class="form-control"
                  aria-describedby="emailHelp"
                  placeholder="Enter UserName"
                  v-model="newcontact"
                />
              </div>
              <b-button v-on:click="addContact()">
                Add
              </b-button>
              <b-button class="m-3" @click="hide">Close</b-button>
              <br />
              <strong class="mr-auto">Go&Vue Chat App</strong>
            </form>
          </div>
        </template>
        <div class="px-3 py-2">
          <p>
            <b-list-group v-for="c in  contacts" :key="c">
              <b-list-group-item button v-on:click="setReciever(c)">{{c}}</b-list-group-item>
            </b-list-group>
          </p>
        </div>
      </b-sidebar>
    </div>
   <div v-if="contacts_selected==true">
      <b-card class="chat-content" title="Chat Page">
        <b-card-text v-for="(m, index) in smsgs" :key="index" :align="m.direct">
          <b-badge variant="primary">
            <img
              src="http://www.bitrebels.com/wp-content/uploads/2011/02/Original-Facebook-Geek-Profile-Avatar-1.jpg"
              class="img-responsive"
            />
          </b-badge>
          <b-card align="left" :sub-title="m.sender" class="chat-msg">
            <b-card-text>{{m.msg}}</b-card-text>
          </b-card>
        </b-card-text>
      </b-card>
      <b-form inline>
        <b-textarea
        rows="0"
        max-rows="1"
        size="sm"
        id="inline-form-input-name"
        v-model="msg"
        class="mb-2 mr-sm-2 mb-sm-0"
        placeholder="...."
        ></b-textarea>
        <b-button v-on:click="sendMessage()" variant="primary">Send</b-button>
        </b-form>
    </div>
    <div  v-else>
      <b-card class="chat-content"></b-card>
    </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import Vue from "vue";
import { BootstrapVue, IconsPlugin } from "bootstrap-vue";

// Install BootstrapVue
Vue.use(BootstrapVue);
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin);
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
export default {
  name: "App",
  data: function () {
    return {
      connection: null,
      msg: "",
      smsgs: [],
      contacts: [],
      menuVisible: false,
      sender: "",
      reciever: "",
      username:"",
      logged:false,
      contacts_selected:false,
      newcontact:""
    };
  },
  methods: {
    sendMessage: function () {
      console.log("Hello from client ");
      console.log("connection info : ", this.connection);
      this.connection.send(
        JSON.stringify({
          sender: this.sender,
          reciver: this.reciever,
          content: this.msg,
        })
      );
      //console.log(this.reciever)
      this.smsgs.push({ sender: this.sender, msg: this.msg, direct: "left" });
      this.msg = "";
      this.connection.onmessage = this.setMessage;

      //this.msg = message
    },
    setMessage(m) {
      var jmsg = JSON.parse(m.data);
      this.smsgs.push({
        sender: jmsg.sender,
        msg: jmsg.content,
        direct: "right",
      });
    },
    setReciever(r) {
      this.reciever = r;
      this.smsgs=[]
      if(this.sender=="" || this.sender==" "){
        this.contacts_selected = false;
        
      }else{
        this.contacts_selected = true;
        
      }
    },
    loggIn(){
      if (this.username == "" || this.username == " "){
        alert("please enter username to start chat")
    }
    else{
      this.connection.send(JSON.stringify({
        sender:this.username,
        reciver:"server",
        content:"genesis message"
      }));
      this.logged = true
      this.sender = this.username
      axios.post('http://localhost:8888/adduser', { username: this.username
      })
      .then(response => console.log(response))
      .catch(e => {
      console.log("err here ocurred:  ",e);
      })
      }
      var url = 'http://localhost:8888/getcontact/'+ this.username
      //var url = 'http://localhost:8888/getcontact/majid'
      //console.log(url)
      axios
      .get(url)
      .then(response => this.contacts=response.data.Contact.split(","). filter(el =>{
        return el != '';
      }))
      .catch(error => console.log(error))
          
      
    },
    addContact(){
      if(this.newcontact != "" & this.newcontact != " "){
      axios.post('http://localhost:8888/addcontact', { 
        username:this.username,
        contact: this.newcontact
      })
      .then(response => console.log(response))
      .catch(e => {
      console.log("err here ocurred:  ",e);
        
      })
    this.contacts.push(this.newcontact)
    this.newcontact = ""
      }
      
    }
    
  },
  created: function () {
    console.log("Starting connection to WebSocket Server");
    this.connection = new WebSocket("ws://localhost:8888/echo");

    this.connection.onmessage = this.setMessage;
    /*this.connection.onmessage = function(event) {
      console.log("message recived from server ", event);
      console.log(event.data)
      this.msg = event.data
    }*/

    this.connection.onopen = function (event) {
      console.log(event);
      console.log("Successfully connected to the echo websocket server...");
    };
  },
};
</script>

<style scoped>
.chat-content {
  background-color: azure;
  overflow: scroll;
  height: 450px;
  max-height: 450px;
}
.img-responsive {
  max-width: 25px;
  max-height: 25px;
}
.chat-msg {
  max-width: 275px;
  max-height: 175;
}
</style>
