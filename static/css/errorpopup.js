document.addEventListener("DOMContentLoaded", function () {
    var popup = document.getElementById("errorPopup");
    var overlay = document.getElementById("popupOverlay");

    if (popup) {
        popup.style.display = "block";
        overlay.style.display = "block";
    }
});

function closePopup() {
    document.getElementById("errorPopup").style.display = "none";
    document.getElementById("popupOverlay").style.display = "none";
}