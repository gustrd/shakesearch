# ShakeSearch

Welcome to ShakeSearch, a web app that allows users to search for a text string in the complete works of Shakespeare. This repository contains the code for Gustavo's Pulley Shakesearch Take-home Challenge.

You can try a live version of the app at https://gustrd-shakesearch.onrender.com/ . Search for "Macbeth" to display a set of results.

The original problem was analyzed from the  **user's perspective**, so the changes do not changed the technology stack,
neither the original architecture. All the focus was at creating new features.

## New Features

I have added several new features to ShakeSearch, including:

- Hygienization of the complete work file to remove references to Project Gutenberg, which previously displayed irrelevant results;
- Case-insensitive search to make it easier for users to find what they are looking for;
- A loading animation to let users know when they are awaiting results;
- User-friendly success and error messages;
- An advanced option for Match Whole Word search;
- Highlighting of the searched word or sentence in yellow;
- A new smart table with pagination to display results more efficiently, including sorting and filtering options;
- Highlighting of filtered word parts, words, or sentences in silver;
- A second column in the results to show which play and act the system identified the resulting text to be in;
- Downloading results to a CSV file;
- An advanced option to configure the average length of the resulting texts;
- An advanced option to provide an OpenAI API key, which can be used to correct misspellings if no results were originally found.

## Technical Enhancements

Although the focus of this work was on the user's perspective, I have made several technical enhancements to improve the codebase, including:

- Code documentation with comments;
- Inclusion of unit testing and end-to-end testing to the code;
- A Dockerfile for easy execution at a Docker container;
- Atomic commits with "conventional commit" tags at git.
- Integration of Bootstrap and jQuery libraries to improve user experience;
- Use of the open-source FancyTable library for advanced table features;
- Update of Go version to 1.18.

## Future Changes Priority

If I had more time and resources to develop ShakeSearch, I would consider the following changes:

- Development of a new front-end with React to improve the design and access to a better table library;
- Use of a PostgreSQL database with GORM and an algorithm to save data on structured tables for better performance, line number results, and Scene identification;
- Testing of other LLM solutions to find one with a better cost-result ratio than OpenAI's Da Vinci for correcting misspellings;
- Use of play or scene as an additional query parameter, with autocomplete to select;
- Refactoring of the code into different files for easier maintenance and better understanding.

