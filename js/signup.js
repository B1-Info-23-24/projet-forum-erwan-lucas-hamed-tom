document.getElementById('signup-form').addEventListener('submit', async function(event) {
    event.preventDefault();

    const username = document.getElementById('pseudoLogin').value;
    const email = document.getElementById('emailLogin').value;
    const password = document.getElementById('passwordLogin').value;
    const password2 = document.getElementById('password2Login').value;

    const messageDiv = document.getElementById('message');

    if (password !== password2) {
        messageDiv.textContent = 'Passwords do not match';
        return;
    }

    const response = await fetch('/api/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password, email })
    });

    const result = await response.json();

    if (response.ok) {
        messageDiv.textContent = 'User registered successfully!';
    } else {
        messageDiv.textContent = `Error: ${result.error}`;
    }
});
