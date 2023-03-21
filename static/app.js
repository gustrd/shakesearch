const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}&s=${data.size}&k=${data.key}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results, data.query);
      });
    });
  },

  updateTable: (results, searchTerm) => {
    const tableBody = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push("<tr><td>" + result.Text + "</td><td>" + result.WorkTitle + "</td></tr>");
    }
    tableBody.innerHTML = rows.join('');

    // Highlight all the table cells containing the search term
    $("td").filter(function() {
      // Use a regular expression to match only the exact word, case insensitive
      var regex = new RegExp("\\b" + searchTerm + "\\b", "i");
      return regex.test($(this).text());
    }).html(function(_, html) {
      // Wrap the matching word in a span with a yellow background
      return html.replace(new RegExp("\\b" + searchTerm + "\\b", "gi"), '<span style="background-color: yellow;">$&</span>');
    });
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
