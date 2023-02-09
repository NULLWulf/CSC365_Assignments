
const input = document.getElementById("autocomplete");
const suggestions = document.getElementById("suggestions");

input.addEventListener("input", async function() {
    // Clear the suggestions list
    suggestions.innerHTML = "";

    // Make the HTTP request to get the values
    try {
        const response = await fetch("https://example.com/values");
        if (!response.ok) {
            throw new Error(response.statusText);
        }
        const values = await response.json();

        // Filter the values based on the user's input
        const inputValue = input.value.toLowerCase();
        const filteredValues = values.filter(value =>
            value.toLowerCase().startsWith(inputValue)
        );

        // Add the filtered values to the suggestions list
        filteredValues.forEach(value => {
            const suggestion = document.createElement("div");
            suggestion.innerHTML = value;
            suggestion.addEventListener("click", function() {
                input.value = value;
                suggestions.innerHTML = "";
            });
            suggestions.appendChild(suggestion);
        });
    } catch (error) {
        console.error(error);
    }
});
