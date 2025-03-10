document.addEventListener("DOMContentLoaded", function () {
    document.querySelectorAll("input[type='checkbox']").forEach(checkbox => {
        checkbox.addEventListener("change", function () {
            let inputContainer = this.parentElement.nextElementSibling;
            if (this.checked) {
                inputContainer.style.display = "block";
            } else {
                inputContainer.style.display = "none";
            }
        });
    });
});
