const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}&s=${data.size}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const tableBody = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push("<tr><td>" + result + "</td></tr>");
    }
    tableBody.innerHTML = rows.join('');
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);

$(document).ready(function() {
  // This function will be called when the HTML document has finished loading

  // Call a function when the table is loaded
  $('#table').fancyTable({
    sortColumn:0,
    pagination: true,
    perPage:5,
    globalSearch:true
  }),

  $("#form").validate({
    rules : {
          query:{
                 required:true,
                 minlength:3
          },
          size:{
                 required:true,
                 min: 50,
                 max: 600
          }                              
    },
    messages:{
          query:{
                 required:"You need to provide the word or sentence to search",
                 minlength:"The word needs to have at least 3 characters"
          },
          size:{
                 required:"A size is needed",
                 min: "The minimum average response is 50",
                 max: "The maximum average response is 600"
          } 
    }
  })

});
