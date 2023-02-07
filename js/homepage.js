const searchButton = document.getElementById("searchButton");
const searchBox = document.getElementById("searchBox");
const label1 = document.getElementById("label1");
const label2 = document.getElementById("label2");

searchButton.addEventListener("click", function () {
    const searchTerm = searchBox.value;
    // fetch(`https://api.example.com/search?q=${searchTerm}`)
    //     .then((response) => response.json())
    //     .then((data) => {
    //         label1.innerHTML = data.label1;
    //         label2.innerHTML = data.label2;
    //     });
    fetch(`/api/v1/business/search?q=${searchTerm}`)
        .then((response) => response.json())
        .then((data) => {
            label1.innerHTML = data.label1;
            label2.innerHTML = data.label2;
        });
});
