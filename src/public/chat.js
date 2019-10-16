

var app = new Vue({
  el: '#app',

  data: {
    ws:null,
    logflow:'<div class="list-group">',
    wsstatus: '[disconnected]',
    msg: '',
    message: {type:'',message:'',from:'',to:''},
    user: { Name: ''},
    users:[],
    toid:'',
    plvalue: ''
  },
  delimiters: ['${', '}'],
  

  mounted: function () {

    this.connectws();
    this.fetchUsers();


  },

  methods: {

    connectws: function(){
            var self = this;
            this.ws = new WebSocket('ws://' + window.location.host + '/ws'+window.location.pathname);
             this.ws.onopen = function() {
                  self.wsstatus='connected';
                };
               this.ws.onclose = function (msg) {
                  self.wsstatus='closed,retry 5s...'
                  setTimeout(function() {
                        self.connectws();
                      }, 5000);
                  return
                };
            this.ws.onmessage= function(e) {
                var msg = JSON.parse(e.data);
                if(msg.type=='whoareyou'){
                     self.wsstatus=msg.message;
                     app.fetchUsers();
                    return
                }
                if(self.logflow.length>20000){
                    self.logflow=self.logflow.substr(15000)
                }

                if(msg.from!=''){
                                self.logflow += '<a href="#" class="list-group-item list-group-item-action list-group-item-success">'
                                +msg.from+':'+msg.message+'</a>';

                }else{
                                self.logflow += '<a href="#" class="list-group-item list-group-item-action list-group-item-info">'+msg.message+'</a>';

                }

                var element = document.getElementById('chat-messages');
                element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
            };
    },
    sendmessage: function(){
      var value = this.msg && this.msg.trim()
      if (!value) {
        return
      }
        message={type:'text',message:this.msg,from:'user',to:this.plvalue}
        this.ws.send(JSON.stringify(message))
        this.msg=''
    },
        setpl: function(name){
            this.plvalue=name;
        },

   fetchUsers: function () {
        command = { command: 'alluser' };
        this.$http.post('/im/command/',command,{emulateJSON: true})
        .then(response => response.json())
        .then(result => {
           Vue.set(this.$data, 'users', result);
            console.log("success in fetch alluser")
                  setTimeout(function() {
                                    app.fetchUsers();
                                  }, 3000);
        })
        .catch(err => {
            console.log(err);
        });
    },


  }
});
