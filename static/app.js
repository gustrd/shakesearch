let firstLoad = true;

const Controller = {
  search: (ev) => {
    ev.preventDefault();

    // Start loading spinner
    document.getElementById('loading-screen').style.display = 'block';
    
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}&s=${data.size}&k=${data.key}&mw=${data.matchWord}`)
      .then(
        (response) => {
          if (response.ok) {
            response.json().then((responseJson) => {          
              Controller.updateTable(responseJson.Results, responseJson.Query, responseJson.MatchWholeWord);
              
              //Show result message
              document.getElementById('result-message').style.display = 'block';
              document.getElementById('result-message').textContent = responseJson.Message;
  
              //Stops loading spinner
              document.getElementById('loading-screen').style.display = 'none';
            });
          }
          else{
              //Show result message
              document.getElementById('result-message').style.display = 'block';
              document.getElementById('result-message').textContent = "There was an error making the search, please try again.";
  
              //Stops loading spinner
              document.getElementById('loading-screen').style.display = 'none';
              
              //Hides the table
              document.getElementById('table').style.display = 'none';
          }
          
        });
  },

  updateTable: (results, searchTerm, useMatchWholeWord) => {
    // Inserts the rows
    const tableBody = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push("<tr><td>" + result.Text + "</td><td>" + result.WorkTitle + "</td></tr>");
    }
    tableBody.innerHTML = rows.join('');

    // Shows the table
    var table = document.getElementById("table");
    table.style.display = "table";

    // Call a function when the table is loaded
    if (!firstLoad){
      // Because of a bug at the library this adjust is needed, or the filter becomes doubled.
      $('th[style="padding:2px;"]').remove();
    }

    $('#table').fancyTable({
      sortColumn:0,
      sortable: true,
      pagination: true,
      perPage:5,
      globalSearch:true,
      inputPlaceholder:"Type here if you want to filter by an additional sentence...",
      onUpdate:function(){
        Controller.atFilter();
      }
    });
      
    firstLoad = false;
    
    // Highlight all the table cells containing the search term
    $("td").filter(function() {
      // Use a regular expression to match the correspondent words
      let regex = new RegExp();
      if (useMatchWholeWord){
        regex = new RegExp("\\b" + searchTerm + "\\b", "i");
      } else {
        regex = new RegExp(searchTerm, "i");
      }

      return regex.test($(this).text());
    }).html(function(_, html) {
      // Wrap the matching word in a span with a yellow background
      return html.replace(new RegExp(searchTerm, "gi"), '<span style="background-color: yellow;">$&</span>');
    });
  },

  atFilter: () => {
    // Add the new color to the filtered words
    var filterTh = $('th[colspan="2"][style="padding:2px;"]').eq(0);
    var filterText = filterTh.find('input').val();
    
    // Cannot let it mess with linebreaks, TODO a better solution, but will need to change the library
    const unfilterableTermsString = "<br><span style=\"background-color: silver\"><span style=\"background-color: yellow\">"
    if (filterText != "" && !unfilterableTermsString.includes(filterText)){
      // Highlight all the table cells containing the search term
      $("td").filter(function() {
        // Use a regular expression to match only the exact word, case insensitive
        var regex = new RegExp(filterText, "i");
        return regex.test($(this).text());
      }).html(function(_, html) {
        // Wrap the matching word in a span with a yellow background
        return html.replace(new RegExp(filterText, "gi"), '<span style="background-color: silver;">$&</span>');
      });
    }
    
  }
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);

$(document).ready(function() {
  // This function will be called when the HTML document has finished loading

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

function toggleDivAdvancedConfigurations() {
  var div = document.getElementById("advanced-configurations");
  if (div.style.display === "none") {
    div.style.display = "block";
  } else {
    div.style.display = "none";
  }
}
