document.addEventListener('DOMContentLoaded', function() {
    var confirmedForm = document.getElementById('confirmed-form');

    if (confirmedForm) {
        confirmedForm.onsubmit = function (onSubmit) {
            onSubmit.preventDefault();

            let title = confirmedForm.getAttribute('data-title');

            if (!title) {
                title = 'Aktion durchführen?';
            }

            return Swal.fire({
                title: title,
                showDenyButton: true,
                confirmButtonText: 'Okay',
                denyButtonText: 'Abbruch',
            }).then((result) => {
                if (result.isConfirmed) {
                    confirmedForm.submit();
                }
            });
        }
    }
});