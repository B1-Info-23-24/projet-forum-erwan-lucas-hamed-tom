document.addEventListener('DOMContentLoaded', function() {
    const passwordInput = document.getElementById('password');
    const toggleIcon = document.querySelector('.password-toggle');
    const eyeOpenIcon = toggleIcon.querySelector('.eye-open');
    const eyeClosedIcon = toggleIcon.querySelector('.eye-closed');

    toggleIcon.addEventListener('click', function() {
        if (passwordInput.type === 'password') {
            passwordInput.type = 'text';
            eyeOpenIcon.style.display = 'block';
            eyeClosedIcon.style.display = 'none';
        } else {
            passwordInput.type = 'password';
            eyeOpenIcon.style.display = 'none';
            eyeClosedIcon.style.display = 'block';
        }
    });
});
