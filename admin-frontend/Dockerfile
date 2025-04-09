# Используем официальный образ Node.js
FROM node:20-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы package.json и package-lock.json (если есть)
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install --frozen-lockfile

# Копируем весь проект в контейнер
COPY . .

# Открываем порт, на котором работает React (по умолчанию 3000)
EXPOSE 5173

# Запускаем приложение
CMD ["npm", "run", "dev", "--", "--host"]