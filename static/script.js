document.addEventListener('DOMContentLoaded', () => {
    fetchTasks();

    document.getElementById('task-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const id = document.getElementById('task-id').value;
        const title = document.getElementById('task-title').value;

        const task = { id: parseInt(id), title, done: false };
        try {
            await createTask(task);
            await fetchTasks(); // Refresh the list after adding
            e.target.reset();
        } catch (error) {
            displayError(`Failed to add task: ${error.message}`);
        }
    });
});

async function fetchTasks() {
    try {
        const response = await fetch('/tasks'); // Use relative path for consistency
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const tasks = await response.json();
        const taskList = document.getElementById('task-list');
        taskList.innerHTML = ''; // Clear the list

        if (tasks.length === 0) {
            taskList.innerHTML = '<li>No tasks available.</li>';
        } else {
            tasks.forEach(task => {
                const li = document.createElement('li');
                li.className = 'task-item';
                li.innerHTML = `
                    <span>${task.title} (ID: ${task.id})</span>
                    <button onclick="deleteTask(${task.id})">Delete</button>
                `;
                taskList.appendChild(li);
            });
        }
    } catch (error) {
        displayError(`Failed to fetch tasks: ${error.message}`);
    }
}

async function createTask(task) {
    const response = await fetch('/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(task),
    });
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
}

async function deleteTask(id) {
    try {
        const response = await fetch(`/tasks/${id}`, { method: 'DELETE' });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        await fetchTasks(); // Refresh the list after deletion
    } catch (error) {
        displayError(`Failed to delete task: ${error.message}`);
    }
}

function displayError(message) {
    const errorDiv = document.getElementById('error-message');
    errorDiv.textContent = message;
    setTimeout(() => errorDiv.textContent = '', 5000); // Clear after 5 seconds
}