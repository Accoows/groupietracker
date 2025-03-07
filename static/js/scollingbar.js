function toggleDropdown() {
    document.getElementById("dropdown").classList.toggle("show");
}

function selectCity(city) {
    // Logique pour sélectionner une ville
    console.log("Ville sélectionnée : " + city);
    // Fermer la barre déroulante après sélection
    document.getElementById("dropdown").classList.remove("show");
    // Réinitialiser le formulaire pour permettre une nouvelle sélection
    document.querySelector('form').reset();
}

window.onclick = function(event) {
    if (!event.target.matches('.dropdown-button')) {
        var dropdowns = document.getElementsByClassName("dropdown-content");
        for (var i = 0; i < dropdowns.length; i++) {
            var openDropdown = dropdowns[i];
            if (openDropdown.classList.contains('show')) {
                openDropdown.classList.remove('show');
            }
        }
    }
}