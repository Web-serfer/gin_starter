// Базовый JavaScript для сайта
document.addEventListener("DOMContentLoaded", function () {
  console.log("Gin приложение загружено!");

  // Анимация появления изображения на главной странице
  const imageElements = document.querySelectorAll('img[src*="optimized-image"]');
  imageElements.forEach((img, index) => {
    // Устанавливаем начальное состояние
    img.style.opacity = "0";
    img.style.transform = "scale(0.8) translateY(20px)";

    // Анимация при загрузке
    setTimeout(() => {
      img.style.transition = "opacity 0.6s ease-out, transform 0.6s ease-out";
      img.style.opacity = "1";
      img.style.transform = "scale(1) translateY(0)";
    }, 300 + (index * 100)); // Задержка для последовательной анимации
  });
});

// Функция для работы со страницей пользователей
function userData() {
	return {
		users: [],
		loading: false,
		showAddForm: false,
		showConfirmationModal: false,
		showSuccessModal: false,
		showErrorModal: false,
		modalUserId: null,
		modalTitle: '',
		modalMessage: '',
		successMessage: '',
		errorMessage: '',
		newUser: {
			name: '',
			email: ''
		},
		addingUser: false,
		showAddSuccessModal: false,
		addSuccessMessage: '',

		fetchUsers() {
			this.loading = true;
			fetch('/api/users')
				.then(response => response.json())
				.then(data => {
					this.users = data;
				})
				.catch(error => {
					console.error('Error fetching users:', error);
					alert('Ошибка при загрузке пользователей');
				})
				.finally(() => {
					this.loading = false;
				});
		},

		addUser() {
			if (!this.newUser.name || !this.newUser.email) {
				// Показываем окно с ошибкой
				this.showErrorModal = true;
				this.errorMessage = 'Пожалуйста, заполните все поля';
				return;
			}

			this.addingUser = true;
			fetch('/api/users', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(this.newUser)
			})
			.then(response => {
				if (response.ok) {
					return response.json();
				}
				throw new Error('Network response was not ok');
			})
			.then(newUser => {
				// Добавляем нового пользователя в конец списка
				this.users.push(newUser);
				// Очищаем форму
				this.newUser = { name: '', email: '' };
				this.showAddForm = false;
				// Показываем окно об успешном добавлении
				this.showAddSuccessModal = true;
				this.addSuccessMessage = 'Пользователь успешно добавлен!';
			})
			.catch(error => {
				console.error('Error adding user:', error);
				// Показываем окно с ошибкой
				this.showErrorModal = true;
				this.errorMessage = 'Ошибка при добавлении пользователя';
			})
			.finally(() => {
				this.addingUser = false;
			});
		},

		closeAddSuccessModal() {
			this.showAddSuccessModal = false;
			this.addSuccessMessage = '';
		},

		showDeleteConfirmation(userId, userName) {
			this.modalUserId = userId;
			this.modalTitle = 'Подтверждение удаления';
			this.modalMessage = `Вы уверены, что хотите удалить пользователя "${userName}"?`;
			this.showConfirmationModal = true;
		},

		cancelDeletion() {
			this.showConfirmationModal = false;
			this.modalUserId = null;
		},

		confirmDeletion() {
			if (!this.modalUserId) {
				return;
			}

			fetch(`/api/users/${this.modalUserId}`, {
				method: 'DELETE'
			})
			.then(response => {
				if (response.ok) {
					// Удаляем пользователя из списка
					this.users = this.users.filter(user => user.id !== this.modalUserId);
					this.showConfirmationModal = false;
					this.modalUserId = null;
					// Показываем окно об успешном удалении
					this.showSuccessModal = true;
					this.successMessage = 'Пользователь успешно удален!';
				} else {
					throw new Error('Network response was not ok');
				}
			})
			.catch(error => {
				console.error('Error deleting user:', error);
				// Показываем окно с ошибкой
				this.showConfirmationModal = false;
				this.showErrorModal = true;
				this.errorMessage = 'Ошибка при удалении пользователя';
			});
		},

		closeSuccessModal() {
			this.showSuccessModal = false;
			this.successMessage = '';
		},

		closeErrorModal() {
			this.showErrorModal = false;
			this.errorMessage = '';
		}
	};
}
