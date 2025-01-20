document.getElementById('emailForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const subject = document.getElementById('subject').value;
    const body = document.getElementById('body').value;
    const attachment = document.getElementById('attachment').files[0];

    const formData = new FormData();
    formData.append('subject', subject);
    formData.append('body', body);
    if (attachment) {
        formData.append('attachment', attachment);
    }

    try {
        const response = await fetch('/api/send-email', {
            method: 'POST',
            body: formData,
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Failed to send email');
        }

        document.getElementById('status').innerHTML = '<p style="color: green;">Email sent successfully!</p>';
    } catch (error) {
        document.getElementById('status').innerHTML = `<p style="color: red;">${error.message}</p>`;
    }
});
