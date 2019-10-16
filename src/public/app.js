// var Vue = require('vue');


// register the grid component
Vue.component('panel-table', {
  template:'\
               <table class="table table-bordered table-striped specialCollapse">\
                 <thead>\
                 <tr>\
                   <th v-for="key in columns"\
                       @click="sortBy(key)"\
                       :class="{ active: sortKey == key }">\
                     {{ key | capitalize }}\
                       <span class="arrow" :class="sortOrders[key] > 0 ? \'asc\' : \'dsc\'">\
                       </span>\
                   </th>\
                 </tr>\
                 </thead>\
                 <tbody>\
                 <tr v-for="entry in filteredData">\
                   <td v-for="key in columns">\
                     {{entry[key]}}\
                   </td>\
                 </tr>\
                 </tbody>\
               </table>',
  props: {
    data: Array,
    columns: Array,
    filterKey: String,
    title: String
  },
  data: function () {
    var sortOrders = {}
    this.columns.forEach(function (key) {
      sortOrders[key] = 1
    })
    return {
      sortKey: '',
      sortOrders: sortOrders
    }
  },
  computed: {
    filteredData: function () {
      var sortKey = this.sortKey
      var filterKey = this.filterKey && this.filterKey.toLowerCase()
      var order = this.sortOrders[sortKey] || 1
      var data = this.data
      if (filterKey) {
        data = data.filter(function (row) {
          return Object.keys(row).some(function (key) {
            return String(row[key]).toLowerCase().indexOf(filterKey) > -1
          })
        })
      }
      if (sortKey) {
        data = data.slice().sort(function (a, b) {
          a = a[sortKey]
          b = b[sortKey]
          return (a === b ? 0 : a > b ? 1 : -1) * order
        })
      }
      return data
    }
  },
  filters: {
    capitalize: function (str) {
      return str.charAt(0).toUpperCase() + str.slice(1)
    }
  },
  methods: {
    sortBy: function (key) {
           this.sortKey = key
           this.sortOrders[key] = this.sortOrders[key] * -1
         },
     removeEntry: function (key) {
       alert(key["Host"])
     }
  }
})

// register the grid component
Vue.component('panel-list', {
  template:'\
        <div>\
        <nav class="navbar navbar-light bg-light justify-content-between">\
            <a class="navbar-brand">{{title}}:{{data.length}}</a>\
            <form class="form-inline">\
                <input class="form-control mr-sm-2" type="search" placeholder="Search" aria-label="Search" v-model="filterKey">\
                <button class="btn btn-danger pull-right" v-on:click="clearBlockmap">Clear Blockmap</button>\
            </form>\
        </nav>\
              <div class="list-group">\
                <a href="#" v-for="(entry,index) in filteredData" class="list-group-item list-group-item-action flex-column align-items-start">\
                    <div class="d-flex w-100 justify-content-between">\
                      <h5 class="mb-1">{{entry["Host"]}}</h5>\
                      <small @click="removeEntry(index,entry)">delete</small>\
                    </div>\
                        <small class="text-success pull-right">{{entry["Reason"]}}</small>\
                </a>\
              </div></div>',
  props: {
    data: Array,
    columns: Array,
    filterKey: String,
    title: String
  },
  data: function () {
    var sortOrders = {}
    this.columns.forEach(function (key) {
      sortOrders[key] = 1
    })
    return {
      sortKey: '',
      sortOrders: sortOrders
    }
  },
  computed: {
    filteredData: function () {
      var sortKey = this.sortKey
      var filterKey = this.filterKey && this.filterKey.toLowerCase()
      var order = this.sortOrders[sortKey] || 1
      var data = this.data
      if (filterKey) {
        data = data.filter(function (row) {
          return Object.keys(row).some(function (key) {
            return String(row[key]).toLowerCase().indexOf(filterKey) > -1
          })
        })
      }
      if (sortKey) {
        data = data.slice().sort(function (a, b) {
          a = a[sortKey]
          b = b[sortKey]
          return (a === b ? 0 : a > b ? 1 : -1) * order
        })
      }
      return data
    }
  },
  filters: {
    capitalize: function (str) {
      return str.charAt(0).toUpperCase() + str.slice(1)
    }
  },
  methods: {
    sortBy: function (key) {
           this.sortKey = key
           this.sortOrders[key] = this.sortOrders[key] * -1
         },
     removeEntry: function (index, entry) {

           console.log("change remove stat")
             command = { command: 'removelblockmap',host:entry["Host"] };
             this.$http.post('/api/command/',command,{emulateJSON: true})
             .then(response => response.json())
             .then(result => {
                this.data.splice(index, 1);
                 console.log("success in change stat")
             })
             .catch(err => {
                 console.log(err);
             });
     },
        clearBlockmap: function () {
             console.log("send fetch servers")
             command = { command: 'clearblockmap' };
             this.$http.post('/api/command/',command,{emulateJSON: true})
             .then(response => response.json())
             .then(result => {

                 app.$emit('reloadblockmap','clearb')
             })
             .catch(err => {
                 console.log(err);
             });
         }

  }
})

var app = new Vue({
  el: '#app',

  data: {
    ws:null,
    logflow:'[log from server...      ',
    wsstatus: '[disconnected]',
    stat: { value: '',file:''},
    server: { name: '', total: '', fail: '' ,rt: ''},
    servers: [],
      users: [],
    //panel table parameters
    title: 'Blockmap list',
    searchQuery: '',
    gridColumns: ['Host', 'Reason'],
    gridData: []

  },
  delimiters: ['${', '}'],
  

  mounted: function () {

    this.connectws();
  },

  methods: {

    connectws: function(){
            var self = this;
            this.ws = new WebSocket('ws://' + window.location.host + '/ws/chat/console');
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

                app.fetchUsers();

                if(self.logflow.length>20000){
                    self.logflow=self.logflow.substr(15000)
                }

                self.logflow += '<div class="chip">'+msg.message+'<br/>';

                var element = document.getElementById('chat-messages');
                element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
            };
    },




      fetchUsers: function () {
          command = { command: 'alluser' };
          this.$http.post('/im/command/',command,{emulateJSON: true})
              .then(response => response.json())
              .then(result => {
                  Vue.set(this.$data, 'users', result);
                  console.log("success in fetch alluser")
              })
              .catch(err => {
                  console.log(err);
              });
      }
  }
});
